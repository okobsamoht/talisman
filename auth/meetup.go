package auth

import (
	"github.com/okobsamoht/tomato/errs"
	"github.com/okobsamoht/tomato/types"
	"github.com/okobsamoht/tomato/utils"
)

type meetup struct{}

func (a meetup) ValidateAuthData(authData types.M, options types.M) error {
	host := "https://api.meetup.com/2/"
	path := "member/self"
	headers := map[string]string{
		"Authorization": "bearer " + utils.S(authData["access_token"]),
	}
	data, err := request(host+path, headers)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Meetup.")
	}
	if data["id"] != nil && utils.S(data["id"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Meetup auth is invalid for this user.")
}
