package u2f

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (u2f *U2F) EnrollHandler(u User, w http.ResponseWriter, r *http.Request) {
	e, err := u2f.Enroll(u)
	resp(w, e, err)
}

func (u2f *U2F) BindHandler(u User, w http.ResponseWriter, r *http.Request) {
	err := u2f.Bind(u, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("bind: %s", err), 500)
		return
	}

	fmt.Fprintf(w, "\"true\"")
}

func (u2f *U2F) SignHandler(u User, w http.ResponseWriter, r *http.Request) {
	e, err := u2f.Sign(u)
	resp(w, e, err)
}

func (u2f *U2F) VerifyHandler(u User, w http.ResponseWriter, r *http.Request) {
	e, err := u2f.Verify(u, r.Body)
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
