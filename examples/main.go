package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/badalex/u2f"
)

func authUser(mu2f *u2f.U2F, r *http.Request) (u2f.User, error) {
	u, err := mu2f.UserList.GetUser("test")
	if err != nil {
		mu2f.UserList.PutUser(u2f.User{User: "test"})
		u, err = mu2f.UserList.GetUser("test")
		if err != nil {
			return u, err
		}
	}

	return u, nil

	// for now we hardcode test, a real example might look like:
	//	user := r.URL.Query()["username"]
	//	if len(user) == 0 {
	//		http.Error(w, fmt.Sprintf("no username passed"), 500)
	//		return ""
	//	}
	//
	//	if user[0] == "" {
	//		http.Error(w, fmt.Sprintf("no username passed"), 500)
	//		return ""
	//	}
	//
	//	// XXX check password
	//	return user[0]
}

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

func (ud *userDB) PutUser(u u2f.User) {
	ud.lock.Lock()
	defer ud.lock.Unlock()

	if ud.Users == nil {
		ud.Users = make(map[string]u2f.User)
	}
	if len(u.User) == 0 {
		panic("wtf")
	}
	ud.Users[u.User] = u
}

func main() {
	var udb = userDB{}
	var mu2f = &u2f.U2F{
		UserList: &udb,
		AppID:    "http://localhost:8081",
		Version:  "U2F_V2",
	}

	http.HandleFunc("/enroll", func(w http.ResponseWriter, r *http.Request) {
		u, err := authUser(mu2f, r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if len(u.User) == 0 {
			panic("WTF")
		}

		mu2f.EnrollHandler(u, w, r)
	})

	http.HandleFunc("/bind", func(w http.ResponseWriter, r *http.Request) {
		u, err := authUser(mu2f, r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		mu2f.BindHandler(u, w, r)
	})

	http.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		u, err := authUser(mu2f, r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		mu2f.SignHandler(u, w, r)
	})

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		u, err := authUser(mu2f, r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		mu2f.SignHandler(u, w, r)
	})

	log.Fatal(http.ListenAndServe(":8079", nil))
}
