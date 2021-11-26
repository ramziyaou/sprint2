package modelz

type Wallet struct {
	ID        int    `json:"id"`
	Ts        string `json:"ts"`
	AccountNo string `json:"accountno"`
	IIN       string `json:"iin"`
	Amount    int    `json:"amount"`
}
