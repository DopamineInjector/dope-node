package communication

import (
	"encoding/json"
	"net/http"
)

func handleAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var input struct {
			PublicKey  string `json:"publicKey"`
			PrivateKey string `json:"privateKey"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// logic
		w.WriteHeader(http.StatusCreated)
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
			Balance   int32  `json:"balance"`
		}

		// logic
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
