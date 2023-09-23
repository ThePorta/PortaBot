package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ThePorta/PortaBot/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
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

	maliciousAddressCh := redisClient.PSub(ctx, "maliciousAddress")

	for {
		select {
		case <-ctx.Done():
			logrus.Info("quit bot")
			return
		case update := <-updates:
			if strings.Compare(update.Message.Text, "/start") == 0 {
				optCode := generateOptCode()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hello %s, Your OPT code is: %s", update.Message.From.FirstName, optCode))
				bot.Send(msg)
			}
		case maliciousAddress := <-maliciousAddressCh:
			checkApprove(ctx, maliciousAddress.Payload)
		}
	}
}

func generateOptCode() string {
	return fmt.Sprintf("%06d", rng.Intn(1000000))
}

func checkApprove(ctx context.Context, maliciousAddress string) {
	accounts, err := redisClient.GetAllAccounts(ctx)
	if err != nil {
		logrus.WithError(err).Error("checkApprove: fail to get all accounts")
		return
	}
}
