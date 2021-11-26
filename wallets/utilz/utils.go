package utilz

import (
	"encoding/json"
	"net/http"
	"sprint2/wallets/modelz"
)

func RespondWithError(w http.ResponseWriter, status int, error modelz.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func ResponseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}
