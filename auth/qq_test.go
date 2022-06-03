package auth

import (
	"fmt"
	"testing"

	"github.com/okobsamoht/talisman/types"
)

func Test_qq_ValidateAuthData(t *testing.T) {
	authData := types.M{
		"access_token": "abc",
		"id":           "123",
	}
	a := yixin{}
	err := a.ValidateAuthData(authData, nil)
	fmt.Println(err)
}
