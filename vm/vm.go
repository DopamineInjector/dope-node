package vm

import (
	"bytes"
	"dope-node/config"
	"os/exec"

	log "github.com/sirupsen/logrus"
)


type RunContractOpts struct {
	// path to the contract binary code on disk
	BinaryPath string;
	// name of the function to be executed
	Entrypoint string;
	// function arguments in serialized string form
	Args string;
	// Address of current block (probably optional)
	BlockNumber string;
	// Transaction sender wallet address
	Sender string;
	// Transaction blockchain id
	TransactionId string;
}

func getVmAddress() string {
	return config.GetString(config.VmAddressKey);
}

func getDbAdress() string {
	return config.GetString(config.DbUrlKey);
}

func RunContract(options *RunContractOpts) (string, error) {
	vmAddress := getVmAddress();
	dbAddress := getDbAdress();
	// Why? I actually ran out of ideas on how to fix this piece of shit system.
	// Do kebab kraftowy like in Kebab Emporium Gdańsk Chełm
	cmd := exec.Command(vmAddress, "-b", options.BinaryPath, "--blockaddr", options.TransactionId, "-d", dbAddress, "-s", options.Sender, "--block-number", options.BlockNumber, "-e", options.Entrypoint, "-a", options.Args)
	var stdout, stderr bytes.Buffer;
	cmd.Stderr = &stderr;
	cmd.Stdout = &stdout;
	err := cmd.Run();
	if err != nil {
		log.Warnf("args")
		for _, a := range(cmd.Args) {
			log.Warn(a)
		}
		return string(stderr.Bytes()), err
	}
	return string(stdout.Bytes()), nil
}
