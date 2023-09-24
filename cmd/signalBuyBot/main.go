package main

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ThePorta/PortaBot/config"
	"github.com/ThePorta/PortaBot/redis"
	"github.com/ThePorta/PortaBot/utils"
	"github.com/ethereum/go-ethereum/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/lmittmann/w3"
	"github.com/sirupsen/logrus"
)

var rng *rand.Rand
var redisClient *redis.Redis

const USDC = "0x07865c6E87B9F70255377e024ace6630C1Eaa37F"
const WETH = "0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6"
const USER = "0xEb4D393215857c51be4e5e61ea8aeabAA7cE9C52"
const UNISWAP_SR02 = "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45"

func init() {
	err := godotenv.Load(".env2")
	if err != nil {
		logrus.WithError(err).Fatal("load .env2")
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logrus.WithError(err).Fatal("set log level")
	}
	logrus.SetLevel(logLevel)
	s := rand.NewSource(time.Now().UnixNano())
	rng = rand.New(s)

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logrus.WithError(err).Fatal("redis db is not a number")
	}
	redisClient = redis.NewRedis(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PWD"), db)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGQUIT)
	go func() {
		<-signalCh
		logrus.Println("waiting for program to quit")
		cancel()
	}()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		logrus.WithError(err).Fatal("init bot")
	}
	bot.Debug = logrus.GetLevel() == logrus.DebugLevel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	logrus.Info("start bot")

	signalBuyCh := redisClient.PSub(ctx, redis.ETH_BUY_SIGNAL)

	var chatId int64
	for {
		select {
		case <-ctx.Done():
			logrus.Info("quit bot")
			return
		case update := <-updates:
			if update.Message != nil && strings.Compare(update.Message.Text, "/start") == 0 {
				chatId = update.Message.Chat.ID
				optCode := generateOptCode()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hello %s, Your OTP code is: %s", update.Message.From.FirstName, optCode))
				bot.Send(msg)
			} else if update.CallbackQuery != nil {
				ans := update.CallbackQuery.Data
				if ans != "no" {
					amount, err := strconv.Atoi(ans)
					if err != nil {
						panic(err)
					}
					decimalBase := big.NewInt(0)
					decimalBase = decimalBase.Exp(big.NewInt(10), big.NewInt(18), nil)
					weiAmount := big.NewInt(int64(amount))
					weiAmount = weiAmount.Mul(weiAmount, decimalBase)
					execBuyETH(ctx, bot, chatId, weiAmount, USER)
				}
			}
		case <-signalBuyCh:
			askBuyETH(ctx, chatId, bot)
		}
	}
}

func generateOptCode() string {
	return fmt.Sprintf("%06d", rng.Intn(1000000))
}

func askBuyETH(ctx context.Context, chatId int64, bot *tgbotapi.BotAPI) {
	buyQuery := tgbotapi.NewMessage(chatId, "ETH is rising, buy some?")
	buyQuery.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("buy 30 USDC", "30"),
			tgbotapi.NewInlineKeyboardButtonData("buy 50 USDC", "50"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("buy 100 USDC", "100"),
			tgbotapi.NewInlineKeyboardButtonData("no", "no"),
		),
	)
	bot.Send(buyQuery)
}
func execBuyETH(ctx context.Context, bot *tgbotapi.BotAPI, chatId int64, amountIn *big.Int, recipient string) {
	// in, minOut, path, recipient
	uniFunc := w3.MustNewFunc("swapExactTokensForTokens(uint256,uint256,address[],address)", "uint256")

	input, err := uniFunc.EncodeArgs(amountIn, big.NewInt(0), []common.Address{common.HexToAddress(USDC), common.HexToAddress(WETH)}, common.HexToAddress(recipient))
	if err != nil {
		panic(err)
	}

	chain := config.ChainConfigs[0]
	// inputHexString := hexutil.Encode(input)
	uuidStr := uuid.New().String()
	err = redisClient.SetInputData(ctx, uuidStr, recipient, utils.Bytes2Hex(input), chain.ChainId, chain.ChainName, UNISWAP_SR02)
	if err != nil {
		panic(err)
	}

	logrus.Info(fmt.Sprintf("UUID : %s", uuidStr))
	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("http://localhost:3000/submitTx?uuid=%s", uuidStr))
	bot.Send(msg)
}
