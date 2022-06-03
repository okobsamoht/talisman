package auth

import (
	"github.com/okobsamoht/talisman/errs"
	"github.com/okobsamoht/talisman/types"
	"github.com/okobsamoht/talisman/utils"
)

type github struct{}

func (a github) ValidateAuthData(authData types.M, options types.M) error {
	host := "https://api.github.com/"
	path := "user"
	headers := map[string]string{
		"Authorization": "bearer " + utils.S(authData["access_token"]),
		"User-Agent":    "parse-server",
	}
	data, err := request(host+path, headers)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Github.")
	}
	if data["id"] != nil && utils.S(data["id"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Github auth is invalid for this user.")
}
