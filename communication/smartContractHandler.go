package communication

import (
	"dope-node/blockchain"
	"dope-node/communication/messages"
	"dope-node/config"
	"dope-node/utils"
	"dope-node/vm"
	"encoding/base64"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func handleSmartContract(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var input messages.SmartContractRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		var output struct {
			Output string `json:"output"`
		}

		marshalledPayload, err := json.Marshal(input.Payload)
		if err != nil {
			log.Warnf("error marshalling payload")
		}
		sndr, err := base64.StdEncoding.DecodeString(input.Payload.Sender)
		if err != nil {
			http.Error(w, "Invalid sender encoding", http.StatusBadRequest)
			return
		}
		sig, err := base64.StdEncoding.DecodeString(input.Signature)
		if err != nil {
			http.Error(w, "Invalid signature encoding", http.StatusBadRequest)
			return
		}
		utils.VerifySignature(sndr, marshalledPayload, sig)

		// Empty output string if not viewing
		output.Output = ""
		if input.View {
			out, err := vm.RunContract(&vm.RunContractOpts{BinaryPath: config.VmAddressKey, Entrypoint: input.Payload.Entrypoint, Args: input.Payload.Args, Sender: string(input.Payload.Sender), TransactionId: blockchain.DopeTransactions[len(blockchain.DopeTransactions)-1].Id})
			if err != nil {
				log.Warnf("error while running VM: %s", err)
			}
			log.Infof("VM output: %s", out)
			output.Output = out
		} else {
			parsedSC := input.ParseToSmartContract()
			blockchain.DopeContracts.SaveContract(&parsedSC)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
