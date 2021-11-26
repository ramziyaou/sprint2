package main

import (
	"log"
	"math/rand"
	"net/http"
	"sprint2/wallets/controllerz"
	"sprint2/wallets/driverz"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	controller := &controllerz.Controller{}
	db, err := driverz.ConnectDB()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/protected", controller.AddWallet(db)).Methods("GET")

	log.Println("Listen on port 8090...")
	log.Fatal(http.ListenAndServe(":8090", router))
}
