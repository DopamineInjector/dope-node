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

type Transactions []Transaction

var DopeTransactions = Transactions{}

func SyncTransactions(transactions *Transactions) {
	DopeTransactions = *transactions
}

func (dTransactions *Transactions) InsertTransaction(transaction *Transaction, dbUrl *string) error {
	if dbUrl == nil || *dbUrl == "" {
		return fmt.Errorf("database URL not set")
	}

	senderBalance, err := utils.GetUserBalance(*dbUrl, transaction.Sender)
	if err != nil {
		return err
	}

	receiverBalance, err := utils.GetUserBalance(*dbUrl, transaction.Receiver)
	if err != nil {
		return err
	}

	newReceiverBalance := receiverBalance + transaction.Amount
	newSenderBalance := senderBalance - transaction.Amount
	if newSenderBalance < 0 {
		return fmt.Errorf("not enouth $")
	}

	_, err = utils.UpddateBalance(*dbUrl, transaction.Sender, newSenderBalance)
	if err != nil {
		return err
	}
	_, err = utils.UpddateBalance(*dbUrl, transaction.Receiver, newReceiverBalance)
	if err != nil {
		return err
	}

	return nil
}

func (dTransactions *Transactions) SaveTransaction(transaction *Transaction) {
	transaction.Id = string(len(*dTransactions))
	*dTransactions = append(*dTransactions, *transaction)
}

func (trans *Transactions) Print() {
	for _, t := range *trans {
		fmt.Println(t)
	}
}

func (t *Transaction) ToString() string {
	return fmt.Sprintf("Transaction: {id: %s, sender: %s, receiver: %s, amount: %d}", t.Id, t.Sender, t.Receiver, t.Amount)
}

// Implementation of transactable
func (t *Transaction) run() (*string, error) {
	dbUrl := config.GetString(config.DbUrlKey);
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
