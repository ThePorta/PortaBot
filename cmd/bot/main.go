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
	"github.com/lmittmann/w3/module/eth"
	"github.com/sirupsen/logrus"
)

var rng *rand.Rand
var redisClient *redis.Redis

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.WithError(err).Fatal("load .env")
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

	maliciousAddressCh := redisClient.PSub(ctx, redis.MALICIOUS_ALTER)

	for {
		select {
		case <-ctx.Done():
			logrus.Info("quit bot")
			return
		case update := <-updates:
			if strings.Compare(update.Message.Text, "/start") == 0 {
				optCode := generateOptCode()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hello %s, Your OTP code is: %s", update.Message.From.FirstName, optCode))
				bot.Send(msg)
				optCodeInt, _ := strconv.Atoi(optCode)
				redisClient.SetOpt2ChatId(ctx, optCodeInt, update.Message.Chat.ID)
			}
		case maliciousAddress := <-maliciousAddressCh:
			checkApprove(ctx, maliciousAddress.Payload, bot)
		}
	}
}

func generateOptCode() string {
	return fmt.Sprintf("%06d", rng.Intn(1000000))
}

func checkApprove(ctx context.Context, maliciousAddress string, bot *tgbotapi.BotAPI) {
	accounts, err := redisClient.GetAllAccounts(ctx)
	if err != nil {
		logrus.WithError(err).Error("checkApprove: fail to get all accounts")
		return
	}

	for _, chain := range config.ChainConfigs {
		client, err := w3.Dial(chain.EndpointUrl)
		if err != nil {
			logrus.WithError(err).Errorf("checkApprove: dial %s", chain.EndpointUrl)
			continue
		}
		defer client.Close()
		funcAllowance := w3.MustNewFunc("allowance(address,address)", "uint256")
		for tokenName, tokenAddr := range chain.SupportedTokens {
			for _, account := range accounts {
				var allowanceAmount big.Int
				client.Call(
					eth.CallFunc(common.HexToAddress(tokenAddr), funcAllowance, common.HexToAddress(account), common.HexToAddress(maliciousAddress)).Returns(&allowanceAmount),
				)
				if allowanceAmount.Cmp(big.NewInt(0)) > 0 {
					chatId, err := redisClient.GetAccountInfo(ctx, account)
					if err != nil {
						logrus.WithError(err).Errorf("checkApprove: get account from redis, account: %s", account)
					}
					uuidStr := uuid.New().String()
					signInputData, err := getApproveInputData(maliciousAddress)
					if err != nil {
						logrus.WithError(err).Error("checkApprove: encode input data")
						continue
					}

					msgConfig := tgbotapi.NewMessage(chatId, fmt.Sprintf("Security Warning: your account %s approve the malicious address %s %s %s on %s. Please open %s/%s to revoke", account, maliciousAddress, allowanceAmount.String(), tokenName, chain.ChainName, config.URL, uuidStr))
					_, err = bot.Send(msgConfig)
					if err != nil {
						logrus.WithError(err).Errorf("checkApprove: send message: chatId: %+v", msgConfig)
						continue
					}
					err = redisClient.SetInputData(ctx, uuidStr, account, utils.Bytes2Hex(signInputData), chain.ChainId, chain.ChainName, tokenAddr)
					if err != nil {
						logrus.WithError(err).Error("checkApprove: set input data")
						continue
					}
				}
			}
		}
	}
}

func getApproveInputData(maliciousAddress string) ([]byte, error) {
	funcApprove := w3.MustNewFunc("approve(address,uint256)", "bool")
	input, err := funcApprove.EncodeArgs(common.HexToAddress(maliciousAddress), big.NewInt(0))
	if err != nil {
		logrus.WithError(err).Errorf("getApproveInputData: encode args: malicious address: %s", maliciousAddress)
		return nil, err
	}
	return input, nil
}
