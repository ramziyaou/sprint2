package models

type User struct {
	ID       int    `json:"id"`
	Ts       string `json:"ts"`
	IIN      string `json:"iin"`
	Username string `json:"username"`
	Password string `json:"password"`
}
