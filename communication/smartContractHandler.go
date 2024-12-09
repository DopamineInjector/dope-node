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
		if input.View {
			stringSender := base64.StdEncoding.EncodeToString(input.Payload.Sender)
			scPath := config.GetString(config.SCAddressKey)
			out, err := vm.RunContract(&vm.RunContractOpts{BinaryPath: scPath, Entrypoint: input.Payload.Entrypoint, Args: input.Payload.Args, Sender: stringSender, TransactionId: "random", BlockNumber: "2137"})

			if err != nil {
				log.Warnf("error while running VM: %s", err)
			}
			log.Infof("VM output: %s", out)
			output.Output = out
		} else {
			if isOk, _ := utils.VerifySignature(input.Payload.Sender, marshalledPayload, input.Signature); !isOk {
				log.Warn("Sent badly signed shit")
				http.Error(w, "Could not verify signature", http.StatusForbidden)
				return
			}
			parsedSC := input.ParseToSmartContract()
			blockchain.DopeTransactables.InsertTransactable(parsedSC)
		}

		if len(blockchain.DopeTransactables) >= MAX_TRANSACTIONS_PER_BLOCK {
			digBlock("bloczek")
		}

		body, _ := json.Marshal(output)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
