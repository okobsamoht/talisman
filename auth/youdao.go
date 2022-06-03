package auth

import (
	"github.com/okobsamoht/tomato/errs"
	"github.com/okobsamoht/tomato/types"
	"github.com/okobsamoht/tomato/utils"
)

type youdao struct{}

func (a youdao) ValidateAuthData(authData types.M, options types.M) error {
	// 具体接口参考：http://note.youdao.com/open/apidoc.html#_Toc304370863
	client := NewOAuth(options)
	client.Host = "http://note.youdao.com"
	client.AuthToken = utils.S(authData["access_token"])
	client.AuthTokenSecret = utils.S(authData["auth_token_secret"])
	data, err := client.Get("/yws/open/user/get.json", nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Youdao.")
	}
	if data["user"] != nil && utils.S(data["user"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Youdao auth is invalid for this user.")
}
