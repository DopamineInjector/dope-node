package blockchain

import (
	"dope-node/communication/messages"
	"fmt"
	"strconv"

	db "github.com/DopamineInjector/go-dope-db"
)

func Transact(transaction messages.Transaction) (*messages.Transaction, error) {
	senderBalance, err := getUserBalance(dbUrl, transaction.Sender)
	if err != nil {
		return nil, err
	}

	receiverBalance, err := getUserBalance(dbUrl, transaction.Receiver)
	if err != nil {
		return nil, err
	}

	newReceiverBalance := receiverBalance + transaction.Amount
	newSenderBalance := senderBalance - transaction.Amount
	if newSenderBalance < 0 {
		return nil, fmt.Errorf("not enouth $")
	}

	_, err = db.InsertValue(dbUrl, prepareInsertValueRequest(transaction.Sender, newSenderBalance))
	if err != nil {
		return nil, err
	}
	_, err = db.InsertValue(dbUrl, prepareInsertValueRequest(transaction.Receiver, newReceiverBalance))
	if err != nil {
		return nil, err
	}

	//addLeafToMPT(&transaction)
	return &transaction, nil
}

func getUserBalance(dbUrl string, user string) (float64, error) {
	balance, err := db.GetValue(dbUrl, db.SelectValueRequest{Key: user, Namespace: "transaction"})
	if err != nil {
		return 0.0, err
	}

	balanceParsed, err := strconv.ParseFloat(balance.Value, 64)
	if err != nil {
		return 0.0, err
	}

	return balanceParsed, nil
}

func prepareInsertValueRequest(key string, value float64) db.InsertValueRequest {
	return db.InsertValueRequest{Key: key, Value: strconv.FormatFloat(value, 'f', 2, 64), Namespace: "transaction"}
}
