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
				Sender     string `json:"sender"`
				Contract   string `json:"contract"`
				Entrypoint string `json:"entrypoint"`
				Args       string `json:"args"`
			} `json:"payload"`
			Signature string `json:"signature"`
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
		isOk, _ := utils.VerifySignature(sndr, marshalledPayload, sig);
		if !isOk {
			log.Warn("Sent badly signed shit")
			http.Error(w, "Could not verify signature", http.StatusForbidden);
			return
		}

		// Empty output string if only viewing
		output.Output = ""
		if input.View {
			scPath := config.GetString(config.SCAddressKey);
			out, err := vm.RunContract(&vm.RunContractOpts{BinaryPath: scPath, Entrypoint: input.Payload.Entrypoint, Args: input.Payload.Args, Sender: string(input.Payload.Sender), TransactionId: "random", BlockNumber: "2137"})
			if err != nil {
				log.Warnf("error while running VM: %s", err)
			}
			log.Infof("VM output: %s", out)
			output.Output = out
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
