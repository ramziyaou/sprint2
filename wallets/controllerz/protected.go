package controllerz

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	// "sprint1/utils"

	// // "sprint1/utils"

	"sprint2/wallets/modelz"
	"sprint2/wallets/utilz"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

type Controller struct{}

// Protected endpoint for token verification
func (c *Controller) AddWallet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// iin, ok := r.Context().Value("IIN").(string)
		// if !ok {
		// 	fmt.Println("no IIN passed")
		// 	return
		// }
		iin := r.FormValue("IIN")
		if iin == "" {
			fmt.Println("no IIN passed")
			return
		}
		var error modelz.Error
		var prevAccountNo string
		if err := db.QueryRow("select `accountno` from wallets ORDER BY id DESC LIMIT 1").Scan(&prevAccountNo); err != nil {
			error.Message = err.Error()
			utilz.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		var wallet modelz.Wallet
		v, ok := GenerateAccountNo(prevAccountNo[3:])
		if !ok {
			error.Message = "Limit exceeded"
			utilz.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		wallet.AccountNo, wallet.IIN = v, iin
		insForm, err := db.Prepare("insert into wallets (iin, accountno) values(?, ?)")
		if err != nil {
			log.Println(err.Error())
			error.Message = "Server error."
			utilz.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		if _, err := insForm.Exec(wallet.IIN, wallet.AccountNo); err != nil {
			error.Message = err.Error()
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				error.Message = "Accountno already exists"
			}
			utilz.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		fmt.Fprintf(w, "Success", wallet)
	}
}

//const Nums = "0123456789"

func GenerateAccountNo(prev string) (string, bool) {
	num, _ := strconv.Atoi(prev)
	num++
	if s := strconv.Itoa(num); len(s) == 10 {
		return "KZT" + s, true
	}
	return "", false
}
