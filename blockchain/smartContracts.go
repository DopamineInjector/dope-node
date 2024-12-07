package blockchain

import (
	"dope-node/config"
	"dope-node/vm"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type SmartContract struct {
	Sender     []byte `json:"sender"`
	Contract   []byte `json:"contract"`
	Entrypoint string `json:"entrypoint"`
	Args       string `json:"args"`
}

// Transactable

func (t SmartContract) run() (*string, error) {
	scPath := config.GetString(config.SCAddressKey)
	out, err := vm.RunContract(&vm.RunContractOpts{BinaryPath: scPath, Entrypoint: t.Entrypoint, Args: t.Args, Sender: string(t.Sender), TransactionId: "rand", BlockNumber: "(ro)blok"})
	if err != nil {
		log.Warnf("error while running VM: %s", err.Error())
		return &out, err
	}
	log.Infof("VM output: %s", out)
	return &out, nil
}

func (t SmartContract) Print() {
	fmt.Printf("SC [Sender: %s, Entrypoint: %s]", t.Sender, t.Entrypoint)
}
