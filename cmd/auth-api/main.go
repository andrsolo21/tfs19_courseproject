package main

import (
	"courseproject/internal/user"
	"courseproject/internal/auth"
	"encoding/json"
	"fmt"
	"time"

	//"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
)


func signup(w http.ResponseWriter, r *http.Request){

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

	//fmt.Fprintf(w, "\nEmail:\n", resp.Email, "\nPassword:\n", resp.Password)
	fmt.Fprintln(w,"id: ", resp.ID)
	//fmt.Println("id: ", resp)
	fmt.Fprintln(w,"fn: ", resp.First_name)
	fmt.Fprintln(w,"ln: ", resp.Last_name)
	fmt.Fprintln(w,"e: ", resp.Email)
	fmt.Fprintln(w,"p: ", resp.Password)
	fmt.Fprintln(w,"birth", resp.Birthday.String())
	w.Write([]byte("\nGreetings 12345!"))
}

func signup2(w http.ResponseWriter, r *http.Request){
	w.Header().Add("X-MY-LOCATION", "ALASKA")
	w.Write([]byte("\nGreetings 12345!"))
}

func signin(w http.ResponseWriter, r *http.Request){

	var resp auth.Logpas

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "err %q\n", err, err.Error())
	} else {
		err = json.Unmarshal(body, &resp)
		if err != nil {
			fmt.Println(w, "can't unmarshal: ", err.Error())
		}
	}

	w.Write([]byte("\nGreetings 12345!"))
}


func main() {

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/api/v1/signup", signup)
		r.Post("/api/v1/signin", signin)
		r.Get("/api/v1/signup", signup2)
	})
	http.ListenAndServe(":5000", r)
}
