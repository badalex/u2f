package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/badalex/u2f"
)

// userDB - dead simple in memory u2f.Users interface implementation
type userDB struct {
	Users map[string]u2f.User
	lock  sync.Mutex
}

func (ud *userDB) GetUser(user string) (u2f.User, error) {
	ud.lock.Lock()
	defer ud.lock.Unlock()

	if ud.Users == nil {
		ud.Users = make(map[string]u2f.User)
		return u2f.User{}, fmt.Errorf("no such user")
	}

	u, ok := ud.Users[user]
	if !ok {
		return u, fmt.Errorf("no such user")
	}
	return u, nil
}

func (ud *userDB) PutUser(u u2f.User) error {
	ud.lock.Lock()
	defer ud.lock.Unlock()

	if ud.Users == nil {
		ud.Users = make(map[string]u2f.User)
	}
	if u.User == "" {
		return fmt.Errorf("missing username")
	}

	ud.Users[u.User] = u
	return nil
}

// normally you would want to do things like authenticate the user/password or
// see if they have a session here
func authUser(s u2f.U2FServer, w http.ResponseWriter, r *http.Request, cb func(u u2f.User)) {
	u, err := s.Users.GetUser("test")

	// don't have that user yet? just add them for now
	if err != nil {
		u.User = "test"
		err = s.Users.PutUser(u)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	cb(u)
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

func main() {
	s := u2f.StdU2FServer(&userDB{}, "https://gou2f.com:8079")

	http.HandleFunc("/Register", func(w http.ResponseWriter, r *http.Request) {
		authUser(s, w, r, func(u u2f.User) {
			e, err := s.Register(u)
			resp(w, e, err)
		})
	})

	http.HandleFunc("/RegisterFin", func(w http.ResponseWriter, r *http.Request) {
		authUser(s, w, r, func(u u2f.User) {
			err := s.RegisterFin(u, r.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("bind: %s", err), 500)
				return
			}

			w.Header().Set("Content-Type", "text/json")
			fmt.Fprintf(w, "{\"ok\": true}")
		})
	})

	http.HandleFunc("/Sign", func(w http.ResponseWriter, r *http.Request) {
		authUser(s, w, r, func(u u2f.User) {
			e, err := s.Sign(u)
			resp(w, e, err)
		})
	})

	http.HandleFunc("/SignFin", func(w http.ResponseWriter, r *http.Request) {
		authUser(s, w, r, func(u u2f.User) {
			e, err := s.SignFin(u, r.Body)
			resp(w, e, err)
		})
	})

	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	fmt.Println("Listening on :8079")
	log.Fatal(http.ListenAndServeTLS(":8079", "cert.pem", "key.pem", nil))
}
