package communication

import (
	"encoding/json"
	"net/http"
)

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var input struct {
			Payload struct {
				Sender    string `json:"sender"`
				Recipient string `json:"recipient"`
				Amount    int32  `json:"amount"`
			} `json:"payload"`
			Signature string `json:"signature"`
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
