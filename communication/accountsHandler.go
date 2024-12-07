package communication

import (
	"dope-node/utils"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func handleAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var input struct {
			PublicKey string `json:"publicKey"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		utils.RegisterAccount(dbUrl, input.PublicKey)
		w.WriteHeader(http.StatusCreated)
		log.Infof("%s created", input.PublicKey)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAccountsInfo(w http.ResponseWriter, r *http.Request) {
	// PUT, bcs in some implmenetations GET doesn't allow to have body ~Tymek
	if r.Method == http.MethodPut {
		var input struct {
			PublicKey string `json:"publicKey"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		var output struct {
			PublicKey string `json:"publicKey"`
			Balance   int    `json:"balance"`
		}

		balance, err := utils.GetUserBalance(dbUrl, input.PublicKey)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			log.Warnf("Error while receiving user's balance: %s", input.PublicKey)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			output.Balance = balance
			output.PublicKey = input.PublicKey
			json.NewEncoder(w).Encode(output)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
