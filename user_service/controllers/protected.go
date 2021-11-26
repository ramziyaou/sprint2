package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Protected endpoint for token verification
func (c *Controller) ProtectedEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("protected invoked")
		iin, ok := r.Context().Value("IIN").(string)
		if !ok {
			fmt.Println("no IIN passed")
			return
		}
		fmt.Println("passing IIN", iin)
		res, err := http.Get("http://localhost:8090/protected?IIN=" + iin)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// req, err := http.NewRequest("GET", "http://localhost:8090/protected", nil)
		// if err != nil {
		// 	log.Fatalf("%v", err)
		// }

		// ctx := context.WithValue(context.Background(), "IIN", "980124450084")

		// req = req.WithContext(ctx)

		// client := http.DefaultClient
		// res, err := client.Do(req)
		// if err != nil {
		// 	log.Fatalf("%v", err)
		// }

		// req, err := http.NewRequest("GET", "http://localhost:8090/protected", nil)
		// if err != nil {
		// 	log.Fatalf("%v", err)
		// }

		// ctx := context.WithValue(req.Context(), "IIN", "980124450084")

		// req = req.WithContext(ctx)

		// client := http.DefaultClient
		// res, err := client.Do(req)
		// if err != nil {
		// 	log.Fatalf("%v", err)
		// }
		responseData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(responseData))
		fmt.Printf("%v\n", res.StatusCode)
		fmt.Fprintf(w, string(responseData))
		// response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")

		// if err != nil {
		// 	fmt.Print(err.Error())
		// 	os.Exit(1)
		// }

		// responseData, err := ioutil.ReadAll(response.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println(string(responseData))
		fmt.Println("protectedEndpoint invoked.")
	}
}
