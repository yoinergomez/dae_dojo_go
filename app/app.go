package main

import (
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
        Name string
        Phone string
				Cedula string
				Role string
}

///Users/usuario/go/src/github.com/gorilla/mux/
func HomeHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/index.html")
}
func getAllPosts(w http.ResponseWriter, r *http.Request){
	resp,errReques := http.Get("https://jsonplaceholder.typicode.com/posts")
	if errReques != nil {
		w.WriteHeader(405)
		w.Write([]byte("error API"))
		return
	}
	responseData, errParse := ioutil.ReadAll(resp.Body)
	if errParse != nil {
		w.WriteHeader(405)
		w.Write([]byte("error parsing"))
		return
	}
	responseString := string(responseData)
	w.WriteHeader(200)
	w.Write([]byte(responseString))
}

func getAll(w http.ResponseWriter, r *http.Request){


	// fmt.Println("r: ", usuario.Name)
	// session, err := mgo.Dial("mongodb://admin:admin@ds115671.mlab.com:15671/dojogo")
	session, err := mgo.Dial("mongodb:27017/dojogo")
	if err != nil {
					panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	var usuarios []Person
	c := session.DB("dojogo").C("people")
	err = c.Find(bson.M{}).All(&usuarios)
	respuesta, err :=  json.Marshal(usuarios)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respuesta)

}

func create(w http.ResponseWriter, r *http.Request){
	var usuario Person
	if r.Body == nil {
		w.WriteHeader(400)
		w.Write([]byte("Ingresar los datos del usuario"))
		return
	}
	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		fmt.Println("Error Decoder")
		w.WriteHeader(400)
		w.Write([]byte("Ingresar los datos del usuario"))
		return
	}

	// fmt.Println("r: ", usuario.Name)
	// session, err := mgo.Dial("mongodb://admin:admin@ds115671.mlab.com:15671/dojogo")
	session, err := mgo.Dial("mongodb:27017/dojogo")
	if err != nil {
					panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("dojogo").C("people")
	err = c.Insert(&Person{usuario.Name, usuario.Phone, usuario.Cedula, usuario.Role})
	if err != nil {
					log.Fatal(err)
	}
	fmt.Println(usuario.Cedula)
	w.WriteHeader(200)
}

func getPost(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	postId := vars["postId"]
	url := "https://jsonplaceholder.typicode.com/posts/"+postId
	resp,errReques := http.Get(url)
	if errReques != nil {
		w.WriteHeader(405)
		w.Write([]byte("error API"))
		return
	}
	responseData, errParse := ioutil.ReadAll(resp.Body)
	if errParse != nil {
		w.WriteHeader(405)
		w.Write([]byte("error parsing"))
		return
	}
	responseString := string(responseData)
    w.WriteHeader(200)
	w.Write([]byte(responseString))
}

func getByName(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["userName"]
	session, err := mgo.Dial("mongodb:27017/dojogo")
	if err != nil {
					panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	var usuarios []Person
	c := session.DB("dojogo").C("people")
	err = c.Find(bson.M{"Name": name}).All(&usuarios)
	respuesta, err :=  json.Marshal(usuarios)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respuesta)
}

func TestHandler(w http.ResponseWriter, r *http.Request){
	// response, _ := json.Marshal(payload)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(200)
    // w.Write([]byte("blabla\n"))
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["testArg"])

}
func main(){
	fmt.Println("start server")
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
    r.HandleFunc("/", HomeHandler).Methods("GET")
    r.HandleFunc("/get-all-posts", getAllPosts).Methods("GET")
    r.HandleFunc("/get-post/{postId}", getPost).Methods("GET")
	r.HandleFunc("/get-all-users", getAll).Methods("GET")
	r.HandleFunc("/user/{userName}", getByName).Methods("GET")
	r.HandleFunc("/create", create).Methods("POST")
    http.Handle("/", r)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
