package main

import (
	"context"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/ThePorta/PortaBot/redis"
	"github.com/ThePorta/PortaBot/types"
	"github.com/ThePorta/PortaBot/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/w3"
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

		logrus.Info("mock data")
		mockInputData, _ := getApproveInputData("0xba17EEb3F0413b76184bA8Ed73067063FbA6e2eB")

		return c.JSON(http.StatusOK, &types.AccountAndInputData{
			AccountAddress: "0xD0a5266b2515c3b575e30cBC0cfC775FA4fC6660",
			TargetContract: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			InputData:      mockInputData,
			ChainId:        1,
			ChainName:      "Ethereum",
		})
		// return err
	}

	return c.JSON(http.StatusOK, aaid)
}

func setChatId(c echo.Context) error {
	setChatIdReq := new(types.SetChatIdRequest)
	c.Bind(setChatIdReq)
	chatId, err := redisClient.GetOpt2ChatId(context.Background(), setChatIdReq.Otp)
	if err != nil {
		logrus.WithError(err).Errorf("setChatId: otp: %d", setChatIdReq.Otp)
		return err
	}
	err = redisClient.SetAccountInfo(context.Background(), setChatIdReq.Address, chatId)
	if err != nil {
		logrus.WithError(err).Errorf("setChatId: address: %s", setChatIdReq.Address)
		return err
	}

	return c.String(http.StatusOK, "")
}

func getApproveInputData(maliciousAddress string) (string, error) {
	funcApprove := w3.MustNewFunc("approve(address,uint256)", "bool")
	input, err := funcApprove.EncodeArgs(common.HexToAddress(maliciousAddress), big.NewInt(0))
	if err != nil {
		logrus.WithError(err).Errorf("getApproveInputData: encode args: malicious address: %s", maliciousAddress)
		return "", err
	}
	return utils.Bytes2Hex(input), nil
}
