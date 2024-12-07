package blockchain

import (
	"dope-node/config"
	"dope-node/utils"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Transaction struct {
	Id       string
	Sender   string
	Receiver string
	Amount   int
}

func (t *Transaction) ToString() string {
	return fmt.Sprintf("Transaction: {id: %s, sender: %s, receiver: %s, amount: %d}", t.Id, t.Sender, t.Receiver, t.Amount)
}

// Implementation of transactable
func (t Transaction) run() (*string, error) {
	dbUrl := config.GetString(config.DbUrlKey)
	senderBalance, err := utils.GetUserBalance(dbUrl, t.Sender)
	if err != nil {
		return nil, err
	}

	receiverBalance, err := utils.GetUserBalance(dbUrl, t.Receiver)
	if err != nil {
		return nil, err
	}

	newReceiverBalance := receiverBalance + t.Amount
	newSenderBalance := senderBalance - t.Amount
	if newSenderBalance < 0 {
		log.Warn("Transaction: tried to send more money than should")
		return nil, fmt.Errorf("not enough $")
	}

	_, err = utils.UpddateBalance(dbUrl, t.Sender, newSenderBalance)
	if err != nil {
		return nil, err
	}
	_, err = utils.UpddateBalance(dbUrl, t.Receiver, newReceiverBalance)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (trans Transaction) Print() {
	fmt.Printf("Transaction [Id: %s, Sender: %s, Receiver: %s, Amount: %d]", trans.Id, trans.Sender, trans.Receiver, trans.Amount)
}
