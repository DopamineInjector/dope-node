package blockchain

import (
	"dope-node/config"
	"dope-node/vm"

	log "github.com/sirupsen/logrus"
)

type SmartContract struct {
	Sender     []byte `json:"sender"`
	Contract   []byte `json:"contract"`
	Entrypoint string `json:"entrypoint"`
	Args       string `json:"args"`
}

type SmartContracts []SmartContract

var DopeContracts SmartContracts

func (dContract *SmartContracts) SaveContract(contract *SmartContract) {
	*dContract = append(*dContract, *contract)
}

func executeContracts() {
	scPath := config.GetString(config.SCAddressKey)
	for _, c := range DopeContracts {
		out, err := vm.RunContract(&vm.RunContractOpts{BinaryPath: scPath, Entrypoint: c.Entrypoint, Args: c.Args, Sender: string(c.Sender), TransactionId: "rand", BlockNumber: "(ro)blok"})
		if err != nil {
			log.Warnf("error while running VM: %s", err)
		}
		log.Infof("VM output: %s", out)
	}

	DopeContracts = DopeContracts[:0]
}

// Transactable

func (t *SmartContract) run() (*string, error) {
	scPath := config.GetString(config.SCAddressKey)
	out, err := vm.RunContract(&vm.RunContractOpts{BinaryPath: scPath, Entrypoint: t.Entrypoint, Args: t.Args, Sender: string(t.Sender), TransactionId: "rand", BlockNumber: "(ro)blok"})
	if err != nil {
		log.Warnf("error while running VM: %s", err.Error())
		return &out, err
	}
	log.Infof("VM output: %s", out)
	return &out, nil
}
