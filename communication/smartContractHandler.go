package communication

import (
	"dope-node/blockchain"
	"dope-node/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func handleSmartContract(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var input struct {
			Payload struct {
				Sender     []byte `json:"sender"`
				Contract   string `json:"contract"`
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
		utils.VerifySignature([]byte(input.Payload.Sender), string(marshalledPayload), input.Signature)

		if input.View {
			cmd := exec.Command("/bin/dopechain-vm", "-e", input.Payload.Entrypoint, "-a", input.Payload.Args, "-s", input.Payload.Sender, "--block-number", fmt.Sprintf("%d", len(blockchain.DopeChain)-1))
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Warnf("error while running VM: %s", err)
			}
			log.Infof("VM output: %s", out)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
