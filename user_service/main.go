package main

import (
	"log"
	"math/rand"
	"net/http"
	"sprint2/user_service/controllers"
	"sprint2/user_service/driver"

	// "sprint1/controllers"
	// "sprint1/driver"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var tpl *template.Template

func init() {
	gotenv.Load()
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	controller := &controllers.Controller{}
	db, err := driver.ConnectDB()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	router := mux.NewRouter()

	router.HandleFunc("/signup", controller.Signup(db)).Methods("GET", "POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("GET", "POST")
	router.HandleFunc("/protected", controller.TokenVerifyMiddleWare(controller.ProtectedEndpoint())).Methods("GET")

	log.Println("Listen on port 8070...")
	log.Fatal(http.ListenAndServe(":8070", router))
}

// func main() {
//     response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")

//     if err != nil {
//         fmt.Print(err.Error())
//         os.Exit(1)
//     }

//     responseData, err := ioutil.ReadAll(response.Body)
//     if err != nil {
//         log.Fatal(err)
//     }
//     fmt.Println(string(responseData))

// }
