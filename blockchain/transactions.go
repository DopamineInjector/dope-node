package blockchain

import (
	"dope-node/utils"
	"fmt"
)

type Transaction struct {
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
	*dTransactions = append(*dTransactions, *transaction)
}

func (trans *Transactions) Print() {
	for _, t := range *trans {
		fmt.Println(t)
	}
}

func (t *Transaction) ToString() string {
	return fmt.Sprintf("Transaction: {sender: %s, receiver: %s, amount: %d}", t.Sender, t.Receiver, t.Amount)
}
