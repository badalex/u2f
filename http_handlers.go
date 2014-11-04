package u2f

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (f U2F) RegisterHandler(u User, w http.ResponseWriter, r *http.Request) {
	e, err := f.Register(u)
	resp(w, e, err)
}

func (f U2F) RegisterFinHandler(u User, w http.ResponseWriter, r *http.Request) {
	err := f.RegisterFin(u, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("bind: %s", err), 500)
		return
	}

	w.Header().Set("Content-Type", "text/json")
	fmt.Fprintf(w, "{\"ok\": true}")
}

func (f U2F) SignHandler(u User, w http.ResponseWriter, r *http.Request) {
	e, err := f.Sign(u)
	resp(w, e, err)
}

func (f U2F) SignFinHandler(u User, w http.ResponseWriter, r *http.Request) {
	e, err := f.SignFin(u, r.Body)
	resp(w, e, err)
}

func resp(w http.ResponseWriter, e interface{}, err error) {
	if err != nil {
		http.Error(w, fmt.Sprintf("bind: %s", err), 500)
		return
	}

	j := json.NewEncoder(w)
	err = j.Encode(e)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to make json: %s", err), 500)
	}
}
