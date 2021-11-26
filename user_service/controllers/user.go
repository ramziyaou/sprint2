package controllers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"sprint2/user_service/models"
	"sprint2/user_service/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidLength   = errors.New("Please provide a 12-digit IIN")
	ErrNotNumeric      = errors.New("IIN must not contain non-numeric characters")
	ErrInvalidUsername = errors.New("Username must only contain Latin letters, numbers and special characters")
	ErrInvalidPassword = errors.New("Password must only contain Latin letters, numbers and at least one special character")
)

type Controller struct{}

func (c *Controller) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("signup_page.html").Funcs(template.FuncMap{}).ParseFiles("templates/signup_page.html")
		if err != nil {
			fmt.Printf("Parse error: %v\n", err)
			return
		}
		switch r.Method {
		case "GET":
			if err = tmpl.Execute(w, nil); err != nil {
				log.Println(err)
				return
			}
		case "POST":
			var user models.User
			var error models.Error

			IIN, login, password := r.FormValue("iin"), r.FormValue("login"), r.FormValue("password")
			num, err := validateIIN(IIN)
			if err != nil {
				error.Message = err.Error()
				utils.RespondWithError(w, http.StatusBadRequest, error)
				return
			}
			username, pass, err := validateCreds(login, password)
			if err != nil {
				error.Message = err.Error()
				utils.RespondWithError(w, http.StatusBadRequest, error)
				return
			}

			user.IIN, user.Username, user.Password = num, username, pass

			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

			if err != nil {
				log.Fatal(err)
			}

			user.Password = string(hash)

			db, err := sql.Open("mysql", os.Getenv("DATA_SOURCE"))

			// if there is an error opening the connection, handle it
			if err != nil {
				log.Print(err.Error())
				return
			}

			if err := db.Ping(); err != nil {
				fmt.Println(err)
				return
			}

			defer db.Close()
			insForm, err := db.Prepare("insert into users (iin, username, password) values(?, ?, ?)")
			if err != nil {
				log.Println(err.Error())
				error.Message = "Server error."
				utils.RespondWithError(w, http.StatusInternalServerError, error)
				return
			}
			if _, err := insForm.Exec(user.IIN, user.Username, user.Password); err != nil {
				error.Message = err.Error()
				var mysqlErr *mysql.MySQLError
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
					error.Message = "Username exists."
				}
				utils.RespondWithError(w, http.StatusBadRequest, error)
				return
			}
			user.Password = ""
			fmt.Fprintf(w, "Success")
			// http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}

}

func (c *Controller) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("login_page.html").Funcs(template.FuncMap{}).ParseFiles("templates/login_page.html")
		if err != nil {
			fmt.Printf("Parse error: %v\n", err)
			return
		}
		switch r.Method {
		case "GET":
			if err = tmpl.Execute(w, nil); err != nil {
				fmt.Printf("Parse error: %s", err.Error())
				return
			}
		case "POST":

			var user models.User
			var jwt models.JWT
			var error models.Error

			login, password := r.FormValue("username"), r.FormValue("password")
			username, pass, err := validateCreds(login, password)
			if err != nil {
				error.Message = err.Error()
				utils.RespondWithError(w, http.StatusBadRequest, error)
				return
			}

			// user.Username, user.Password = username, pass
			row := db.QueryRow("select * from users where username=?", username)
			err = row.Scan(&user.ID, &user.Ts, &user.IIN, &user.Username, &user.Password)

			if err != nil {
				if err == sql.ErrNoRows {
					error.Message = "The user does not exist"
					utils.RespondWithError(w, http.StatusBadRequest, error)

					return
				} else {
					log.Println(err)
					return
				}
			}
			ctx := context.WithValue(context.Background(), "IIN", user.IIN)

			*r = *r.WithContext(ctx)
			hashedPassword := user.Password

			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pass))

			if err != nil {
				error.Message = "Invalid Password"
				utils.RespondWithError(w, http.StatusUnauthorized, error)
				return
			}

			token, err := utils.GenerateToken(user)

			if err != nil {
				log.Fatal(err)
			}

			w.WriteHeader(http.StatusOK)
			jwt.Token = token

			utils.ResponseJSON(w, jwt)
		}
	}
}

// validateIIN validates IIN
func validateIIN(iin string) (string, error) {
	if len(iin) != 12 {
		return "", ErrInvalidLength
	}

	for _, char := range iin {
		if char < '0' || char > '9' {
			return "", ErrNotNumeric
		}
	}
	return iin, nil
}

// validateCreds validates login and password
func validateCreds(login, password string) (string, string, error) {
	login = strings.Trim(login, " ")
	if !strIsPrint(login) {
		return "", "", ErrInvalidUsername
	}
	if !strIsPrint(password) || !containsSpecialChar(password) {
		return "", "", ErrInvalidPassword
	}
	return login, password, nil
}

// strIsPrint reports whether the string passed consists of printable Latin character only
func strIsPrint(s string) bool {
	for _, char := range s {
		if char < 32 || char > 126 {
			return false
		}
	}
	return true
}

// containsSpecialChar reports whether the string contains a special character
func containsSpecialChar(s string) bool {
	for _, char := range s {
		if char < 48 || (char > 57 && char < 65) || (char > 90 && char < 97) || (char > 122 && char < 127) {
			return true
		}
	}
	return false
}

// TokenVerifyMiddleware verifies token
func (c *Controller) TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject models.Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				return []byte(os.Getenv("SECRET")), nil
			})

			if error != nil {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}

			if token.Valid {
				ctx := context.WithValue(r.Context(), "IIN", "980124450084")

				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Invalid token."
			utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}
