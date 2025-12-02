package middelware

import (
	"fmt"

	"github.com/matthewhartstonge/argon2"
)

func VerifPassword(inputPassword string, hashPassword string) bool {
	ok, err := argon2.VerifyEncoded([]byte(hashPassword),[]byte(inputPassword))

	if err != nil {
		fmt.Println("Error : Password not metmatch, ",err)
	}

	return ok
}