package communication

import (
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
		var input struct {
			Payload struct {
				Sender     []byte `json:"sender"`
				Contract   []byte `json:"contract"`
				Entrypoint string `json:"entrypoint"`
				Args       string `json:"args"`
			} `json:"payload"`
			Signature []byte `json:"signature"`
			View      bool   `json:"view"`
		}
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
		if !input.View {
			if isOk, _ := utils.VerifySignature(input.Payload.Sender, marshalledPayload, input.Signature); !isOk {
				log.Warn("Sent badly signed shit")
				http.Error(w, "Could not verify signature", http.StatusForbidden);
				return
			}
		}
		stringSender := base64.StdEncoding.EncodeToString(input.Payload.Sender)
		scPath := config.GetString(config.SCAddressKey);
		out, err := vm.RunContract(&vm.RunContractOpts{BinaryPath: scPath, Entrypoint: input.Payload.Entrypoint, Args: input.Payload.Args, Sender: stringSender, TransactionId: "random", BlockNumber: "2137"})
		if err != nil {
			log.Warnf("error while running VM: %s", err)
		}
		log.Infof("VM output: %s", out)
		output.Output = out

		body, _ := json.Marshal(output);
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
