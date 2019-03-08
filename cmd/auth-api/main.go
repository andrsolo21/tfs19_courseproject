package main

import (
	"courseproject/internal/user"
	"encoding/json"
	"fmt"
	//"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
)

func GetGreeting(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, anonymous!"))
}
func PostGreeting(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Greetings broadcasted!"))
}


func CreateUserG(w http.ResponseWriter, r *http.Request){


	w.Write([]byte("Greetings 123!"))
}


func CreateUser(w http.ResponseWriter, r *http.Request){

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
	//fmt.Fprintf(w, "\nEmail:\n", resp.Email, "\nPassword:\n", resp.Password)
	fmt.Fprintln(w,"id: ", resp.ID)
	fmt.Fprintln(w,"fn: ", resp.First_name)
	fmt.Fprintln(w,"ln: ", resp.Last_name)
	fmt.Fprintln(w,"e: ", resp.Email)
	fmt.Fprintln(w,"p: ", resp.Password)
	w.Write([]byte("\nGreetings 12345!"))
}


func main() {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/greeting", GetGreeting)
		r.Post("/greeting", PostGreeting)
		r.Get("/add", CreateUserG)
		r.Post("/add", CreateUser)

	})
	http.ListenAndServe(":5000", r)
}

//{ "email": "qweqwe", "password": "3we5"}
//$.post("http://yourserver.com/index.php", {data: { "email": "qweqwe", "password": "3we5"}});
/*data := []byte(`{"foo":"bar"}`)
r := bytes.NewReader(data)
resp, err := http.Post("http://example.com/upload", "application/json", r)
if err != nil {
return err
}*/

//curl --header "Content-Type: applacation/json"\ --request POST\ --data '{ "email": "qweqwe", "password": "3we5"}'\http://localhost:5000/add
//curl --header "Content-Type: applacation/json" --request POST --data '{"id":5, "first_name":"fn","last_name":"ln",  "email":"qwe4qwe", "password":"3we5"}' http: //localhost:5000/add