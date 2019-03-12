package main

import (
	"courseproject/internal/user"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func signup(w http.ResponseWriter, r *http.Request) {

	var resp user.User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "err %q\n", err, err.Error())
	} else {
		err = json.Unmarshal(body, &resp)
		if err != nil {
			fmt.Println(w, "can't unmarshal: ", err.Error())
		}
	}
	resp.Created_at = time.Now()
	resp.Updated_at = time.Now()

	var flag bool
	data, flag = data.AddUser(resp)
	if flag == false {
		fmt.Println("not")
	}
	fmt.Fprintln(w, "otv", flag)

	fmt.Println(data.LenA())
}

func signin(w http.ResponseWriter, r *http.Request) {

	pas := r.PostFormValue("pas")
	log := r.PostFormValue("log")



	fmt.Fprintln(w, "pas: ", pas)
	fmt.Fprintln(w, "log: ", log)
}
