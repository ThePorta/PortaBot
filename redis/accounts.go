package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/ThePorta/PortaBot/types"
	"github.com/sirupsen/logrus"
)

func (r *Redis) SetAccountInfo(ctx context.Context, accountAddress string, chatId int64) (err error) {
	return r.setAndCheck(ctx, accountInfoKey(accountAddress), chatId, "StoreAccountInfo: set")
}

func (r *Redis) GetAccountInfo(ctx context.Context, accountAddress string) (chatId int64, err error) {
	cmd := r.redis.Get(ctx, accountInfoKey(accountAddress))
	chatId, err = cmd.Int64()
	if err != nil {
		logrus.WithError(err).Error("GetAccountInfo: redis get")
		return
	}
	return
}

func (r *Redis) GetAllAccounts(ctx context.Context) (accounts []string, err error) {
	accounts = make([]string, 0)
	iter := r.redis.Scan(ctx, 0, fmt.Sprintf("%s.*", ACCOUNT_INFO), 0).Iterator()
	for iter.Next(ctx) {
		account := strings.Split(iter.Val(), ".")[1]
		accounts = append(accounts, account)
	}
	if err = iter.Err(); err != nil {
		logrus.WithError(err).Error("GetAllAccounts: iter error")
		return []string{}, err
	}
	return
}

func (r *Redis) SetInputData(ctx context.Context, uuid string, accountAddress string, inputData string, chainId int, chainName string, targetContract string) (err error) {
	accountAndInputData := types.AccountAndInputData{
		AccountAddress: accountAddress,
		TargetContract: targetContract,
		InputData:      inputData,
		ChainId:        chainId,
		ChainName:      chainName,
	}
	aaid, err := accountAndInputData.MarshalMsg(nil)
	if err != nil {
		logrus.WithError(err).Errorf("SetInputData: marshal msg: AccountAddress: %s, InputData: %s", accountAddress, inputData)
		return
	}
	return r.setAndCheck(ctx, uuidKey(uuid), aaid, "SetInputData: set")
}

func (r *Redis) GetInputData(ctx context.Context, uuid string) (aaid *types.AccountAndInputData, err error) {
	cmd := r.redis.Get(ctx, uuidKey(uuid))
	aaidBytes, err := cmd.Bytes()
	if err != nil {
		logrus.WithError(err).Error("GetInputData: redis get")
		return
	}
	aaid = new(types.AccountAndInputData)
	_, err = aaid.UnmarshalMsg(aaidBytes)
	if err != nil {
		logrus.WithError(err).Errorf("GetInputData: unmarshal msg: %s", aaidBytes)
		return
	}
	return
}

func (r *Redis) SetOpt2ChatId(ctx context.Context, otp string, chatId int64) error {
	return r.setAndCheck(ctx, opt2ChatIdKey(otp), chatId, "SetOpt2ChatId")
}

func (r *Redis) GetOpt2ChatId(ctx context.Context, otp string) (chatId int64, err error) {
	cmd := r.redis.Get(ctx, opt2ChatIdKey(otp))
	chatId, err = cmd.Int64()
	if err != nil {
		logrus.WithError(err).Error("GetOpt2Chatid: redis get")
		return
	}
	return chatId, nil
}

func accountInfoKey(accountAddress string) string {
	return fmt.Sprintf("%s.%s", ACCOUNT_INFO, strings.ToLower(accountAddress))
}

func uuidKey(uuid string) string {
	return fmt.Sprintf("%s.%s", UUID, uuid)
}

func opt2ChatIdKey(otp string) string {
	return fmt.Sprintf("%s.%s", OTP_2_CHAT_ID, otp)
}
