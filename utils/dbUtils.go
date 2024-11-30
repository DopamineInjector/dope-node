package utils

import (
	"strconv"

	db "github.com/DopamineInjector/go-dope-db"
	log "github.com/sirupsen/logrus"
)

func RegisterAccount(dbUrl string, pub string, priv string) {
	log.Info("New user registered - received 500$ to spend in the casino :)")
	db.InsertValue(dbUrl, db.InsertValueRequest{Key: pub, Value: priv, Namespace: "keys"})
	db.InsertValue(dbUrl, db.InsertValueRequest{Key: pub, Value: "500", Namespace: "balance"})
}

func GetUserBalance(dbUrl string, pub string) (int, error) {
	balance, err := db.GetValue(dbUrl, db.SelectValueRequest{Key: pub, Namespace: "balance"})
	if err != nil {
		return 0, err
	}

	parsedBalance, err := strconv.Atoi(balance.Value)
	if err != nil {
		return 0, err
	}

	return parsedBalance, nil
}

func UpddateBalance(dbUrl string, pub string, value int) (*db.ChecksumResponse, error) {
	return db.InsertValue(dbUrl, prepareInsertValueRequest(&pub, &value))
}

func prepareInsertValueRequest(key *string, value *int) db.InsertValueRequest {
	return db.InsertValueRequest{Key: *key, Value: strconv.Itoa(*value), Namespace: "balance"}
}
