package communication

import (
	"encoding/json"
	"net/http"
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

		// logic
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
