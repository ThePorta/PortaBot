package main

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/ThePorta/PortaBot/redis"
	"github.com/ThePorta/PortaBot/types"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

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
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logrus.WithError(err).Fatal("redis db is not a number")
	}
	redisClient = redis.NewRedis(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PWD"), db)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/getInputData", getInputData)
	e.POST("/setChatId", setChatId)

	e.Logger.Fatal(e.Start(":1234"))
}

func getInputData(c echo.Context) error {
	aaid, err := redisClient.GetInputData(context.Background(), c.QueryParam("uuid"))
	if err != nil {
		logrus.WithError(err).Infof("getInputData: redis get: uuid: %s", c.QueryParam("uuid"))
		return err
	}

	return c.JSON(http.StatusOK, aaid)
}

func setChatId(c echo.Context) error {
	setChatIdReq := new(types.SetChatIdRequest)
	c.Bind(setChatIdReq)
	logrus.Infof("%+v", setChatIdReq)
	chatId, err := redisClient.GetOpt2ChatId(context.Background(), setChatIdReq.Otp)
	if err != nil {
		logrus.WithError(err).Errorf("setChatId: otp: %s", setChatIdReq.Otp)
		return err
	}
	err = redisClient.SetAccountInfo(context.Background(), setChatIdReq.Address, chatId)
	if err != nil {
		logrus.WithError(err).Errorf("setChatId: address: %s", setChatIdReq.Address)
		return err
	}

	return c.String(http.StatusOK, "")
}
