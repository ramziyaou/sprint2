package main

import (
	"fmt"
	"strconv"
)

// func main() {
// 	fmt.Println(validateIIN("980124450084"))
// 	fmt.Println(validateIIN("980124450072"))
// 	fmt.Println(validateIIN("601119400567"))
// 	fmt.Println(validateIIN("910815450350"))

// }

func validateIIN(s string) bool {
	if len(s) != 12 {
		fmt.Println("invalid length!")
		return false
	}
	if _, err := strconv.Atoi(s); err != nil || s[0] == '-' || s[0] == '+' {
		fmt.Println("non-numeric characters!")
		return false
	}
	year := 10*int(s[0]-'0') + int(s[1]-'0')
	//if s[0] > '0' && s[0] < '9' {
	//	fmt.Println("0,9 invalid year!")
	//	return false
	//}
	var month int
	if month := 10*int(s[2]-'0') + int(s[3]-'0'); month > 12 {
		fmt.Println("invalid month!")
		return false
	}
	day := 10*int(s[4]-'0') + int(s[5]-'0')
	if (month == 4 || month == 6 || month == 9 || month == 11) && day > 30 {
		fmt.Println("day shnt exceed 30 for this month!")
		return false
	} else if month == 2 {
		if !isLeapYear(year) {
			if day > 29 {
				fmt.Println(">29 for leap Feb")
				return false
			}
		} else {
			if day > 28 {
				fmt.Println("day > 28 for not leap Feb")
				return false
			}
		}
	} else {
		if day > 31 {
			fmt.Println(">31")
			return false
		}
	}
	if s[6] == '0' || s[6] > '6' {
		fmt.Println("Invalid 7th char")
		return false
	}
	var mod int
	if mod = (int(s[0]-'0') + 2*int(s[1]-'0') + 3*int(s[2]-'0') + 4*int(s[3]-'0') + 5*int(s[4]-'0') + 6*int(s[5]-'0') + 7*int(s[6]-'0') + 8*int(s[7]-'0') + 9*int(s[8]-'0') + 10*int(s[9]-'0') + 11*int(s[10]-'0')) % 11; mod == 10 {
		fmt.Print(mod, "-->")
		mod = (3*int(s[0]-'0') + 4*int(s[1]-'0') + 5*int(s[2]-'0') + 6*int(s[3]-'0') + 7*int(s[4]-'0') + 8*int(s[5]-'0') + 9*int(s[6]-'0') + 10*int(s[7]-'0') + 11*int(s[8]-'0') + int(s[9]-'0') + 2*int(s[10]-'0')) % 11
	}
	fmt.Println(mod)
	if mod != int(s[11]-'0') {
		fmt.Println("invalid last character!")
		return false
	}
	return true
}

func isLeapYear(y int) bool {
	if y%4 == 0 && y%100 != 0 || y%400 == 0 {
		return true
	} else {
		return false
	}
}
