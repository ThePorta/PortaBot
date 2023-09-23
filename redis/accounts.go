package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/ThePorta/PortaBot/types"
	"github.com/sirupsen/logrus"
)

func (r *Redis) StoreAccountInfo(ctx context.Context, accountAddress string, info *types.AccountInfo) (err error) {
	infoDataByte, err := info.MarshalMsg(nil)
	if err != nil {
		logrus.WithError(err).Error("StoreAccountInfo: info marshal msg")
		return
	}
	return r.setAndCheck(ctx, accountInfoKey(accountAddress), infoDataByte, "StoreAccountInfo: set")
}

func (r *Redis) GetAccountInfo(ctx context.Context, accountAddress string) (info *types.AccountInfo, err error) {
	cmd := r.redis.Get(ctx, accountInfoKey(accountAddress))
	infoDatabyte, err := cmd.Bytes()
	if err != nil {
		logrus.WithError(err).Error("GetAccountInfo: redis get")
		return
	}
	_, err = info.UnmarshalMsg(infoDatabyte)
	if err != nil {
		logrus.WithError(err).Error("GetAccountInfo: unmaishaMsg")
	}
	return
}

func (r *Redis) GetAllAccounts(ctx context.Context) (accounts []string, err error) {
	accounts = make([]string, 0)
	iter := r.redis.Scan(ctx, 0, fmt.Sprintf("%s.*", ACCOUNT_INFO), 0).Iterator()
	for iter.Next(ctx) {
		account := strings.TrimSuffix(iter.Val(), fmt.Sprintf("%s.", ACCOUNT_INFO))
		accounts = append(accounts, account)
	}
	if err = iter.Err(); err != nil {
		logrus.WithError(err).Error("GetAllAccounts: iter error")
		return []string{}, err
	}
	return
}

func accountInfoKey(accountAddress string) string {
	return fmt.Sprintf("%s.%s", ACCOUNT_INFO, strings.ToLower(accountAddress))
}
