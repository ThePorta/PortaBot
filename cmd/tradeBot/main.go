package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var rng *rand.Rand

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.WithError(err).Fatal("failed loading .env")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	s := rand.NewSource(time.Now().UnixNano())
	rng = rand.New(s)
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
		logrus.WithError(err).Fatal("TG Token not found")
	}
	bot.Debug = logrus.GetLevel() == logrus.DebugLevel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	logrus.Info("start bot")
	startMonitorTrade(ctx, updates, bot)
}

func startMonitorTrade(ctx context.Context, updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
	for {
		select {
		case <-ctx.Done():
			logrus.Info("quit bot")
			return
		case update := <-updates:
			if update.Message != nil && strings.Compare(update.Message.Text, "/start") == 0 {
				optCode := generateOptCode()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hello %s, Your OPT code is: %s", update.Message.From.FirstName, optCode))
				logrus.Info("send out OTP")
				bot.Send(msg)

				buyQuery := tgbotapi.NewMessage(update.Message.Chat.ID, "Do you wannt buy ETH?")
				buyQuery.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("yes", "yes"),
						tgbotapi.NewInlineKeyboardButtonData("no", "no"),
					),
				)
				bot.Send(buyQuery)
			} else if update.CallbackQuery != nil {
				ans := update.CallbackQuery.Data
				if ans == "yes" {
					// send tx request via WC api (see WC's format)
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Forwarding buy tx via wallet connect"))
				}
			}
		}
	}
}

func generateOptCode() string {
	return fmt.Sprintf("%06d", rng.Intn(1000000))
}
