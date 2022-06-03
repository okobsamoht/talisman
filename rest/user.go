package rest

import (
	"errors"
	"net/url"
	"time"

	"strings"

	"github.com/okobsamoht/talisman/config"
	"github.com/okobsamoht/talisman/errs"
	"github.com/okobsamoht/talisman/mail"
	"github.com/okobsamoht/talisman/orm"
	"github.com/okobsamoht/talisman/types"
	"github.com/okobsamoht/talisman/utils"
)

var adapter mail.Adapter

func init() {
	a := config.TConfig.MailAdapter
	if a == "smtp" {
		adapter = mail.NewSMTPAdapter()
	} else {
		adapter = mail.NewSMTPAdapter()
	}
}

// shouldVerifyEmails 根据配置参数确定是否需要验证邮箱
func shouldVerifyEmails() bool {
	return config.TConfig.VerifyUserEmails
}

// SetEmailVerifyToken 设置需要验证的 token
func SetEmailVerifyToken(user types.M) {
	if user == nil {
		return
	}
	if shouldVerifyEmails() {
		user["_email_verify_token"] = utils.CreateToken()
		user["emailVerified"] = false

		if config.TConfig.EmailVerifyTokenValidityDuration > 0 {
			user["_email_verify_token_expires_at"] = utils.TimetoString(config.GenerateEmailVerifyTokenExpiresAt())
		}
	}
}

// SendVerificationEmail 发送验证邮件
func SendVerificationEmail(user types.M) {
	if shouldVerifyEmails() == false {
		return
	}
	token := url.QueryEscape(utils.S(user["_email_verify_token"]))
	user = getUserIfNeeded(user)
	if user == nil {
		return
	}
	user["className"] = "_User"
	username := url.QueryEscape(utils.S(user["username"]))
	link := buildEmailLink(config.VerifyEmailURL(), username, token)
	options := types.M{
		"appName": config.TConfig.AppName,
		"link":    link,
		"user":    user,
	}
	adapter.SendMail(defaultVerificationEmail(options))
}

// ResendVerificationEmail 重新发送验证邮件
func ResendVerificationEmail(username string) error {
	aUser := getUserIfNeeded(types.M{"username": username})
	if aUser == nil {
		return errors.New("no user")
	}
	if emailVerified, ok := aUser["emailVerified"].(bool); ok && emailVerified {
		return errors.New("emailVerified")
	}
	SetEmailVerifyToken(aUser)
	_, err := orm.TalismanDBController.Update("_User", types.M{"username": username}, aUser, types.M{}, false)
	if err != nil {
		return err
	}
	SendVerificationEmail(aUser)
	return nil
}

// getUserIfNeeded 把 user 填充完整，如果无法完成则返回 nil
func getUserIfNeeded(user types.M) types.M {
	if user == nil {
		return nil
	}
	if user["username"] != nil && user["email"] != nil {
		return user
	}
	where := types.M{}
	if user["username"] != nil {
		where["username"] = user["username"]
	}
	if user["email"] != nil {
		where["email"] = user["email"]
	}

	query, err := NewQuery(Master(), "_User", where, types.M{}, nil)
	if err != nil {
		return nil
	}
	response, err := query.Execute()
	if err != nil {
		return nil
	}
	if utils.HasResults(response) == false {
		return nil
	}
	results := utils.A(response["results"])
	if len(results) != 1 {
		return nil
	}

	return utils.M(results[0])
}

func defaultVerificationEmail(options types.M) types.M {
	if options == nil {
		return nil
	}
	user := utils.M(options["user"])
	if user == nil {
		return nil
	}
	text := "Hi,\n\n"
	text += "You are being asked to confirm the e-mail address " + utils.S(user["email"])
	text += " with " + utils.S(options["appName"]) + "\n\n"
	text += "Click here to confirm it:\n" + utils.S(options["link"])
	to := utils.S(user["email"])
	subject := "Please verify your e-mail for " + utils.S(options["appName"])
	return types.M{
		"text":    text,
		"to":      to,
		"subject": subject,
	}
}

// SendPasswordResetEmail 发送密码重置邮件
func SendPasswordResetEmail(email string) error {
	user := setPasswordResetToken(email)
	if user == nil || len(user) == 0 {
		return errs.E(errs.EmailMissing, "you must provide an email")
	}
	user["className"] = "_User"
	token := url.QueryEscape(utils.S(user["_perishable_token"]))
	username := url.QueryEscape(utils.S(user["username"]))
	link := buildEmailLink(config.RequestResetPasswordURL(), username, token)
	options := types.M{
		"appName": config.TConfig.AppName,
		"link":    link,
		"user":    user,
	}
	adapter.SendMail(defaultResetPasswordEmail(options))
	return nil
}

// setPasswordResetToken 设置修改密码 token
func setPasswordResetToken(email string) types.M {
	token := utils.CreateToken()
	db := orm.TalismanDBController
	where := types.M{
		"$or": types.S{
			types.M{
				"email": email,
			},
			types.M{
				"username": email,
				"email": types.M{
					"$exists": false,
				},
			},
		},
	}
	update := types.M{
		"_perishable_token": token,
	}
	// 增加 token 过期时间
	if config.TConfig.PasswordPolicy && config.TConfig.ResetTokenValidityDuration > 0 {
		update["_perishable_token_expires_at"] = utils.TimetoString(config.GeneratePasswordResetTokenExpiresAt())
	}
	r, err := db.Update("_User", where, update, types.M{}, true)
	if err != nil {
		return nil
	}
	return r
}

func defaultResetPasswordEmail(options types.M) types.M {
	if options == nil {
		return nil
	}
	user := utils.M(options["user"])
	if user == nil {
		return nil
	}
	text := "Hi,\n\n"
	text += "You requested to reset your password for " + utils.S(options["appName"]) + "\n\n"
	text += "Click here to reset it:\n" + utils.S(options["link"])
	var to string
	if utils.S(user["email"]) != "" {
		to = utils.S(user["email"])
	} else {
		to = utils.S(user["username"])
	}
	subject := "Password Reset for " + utils.S(options["appName"])
	return types.M{
		"text":    text,
		"to":      to,
		"subject": subject,
	}
}

// VerifyEmail 更新邮箱验证标志
func VerifyEmail(username, token string) bool {
	if shouldVerifyEmails() == false {
		return false
	}

	db := orm.TalismanDBController
	query := types.M{
		"username":            username,
		"_email_verify_token": token,
	}
	updateFields := types.M{
		"emailVerified": true,
		"_email_verify_token": types.M{
			"__op": "Delete",
		},
	}

	if config.TConfig.EmailVerifyTokenValidityDuration > 0 {
		query["emailVerified"] = false
		query["_email_verify_token_expires_at"] = types.M{
			"$gt": utils.TimetoString(time.Now().UTC()),
		}
		updateFields["_email_verify_token_expires_at"] = types.M{
			"__op": "Delete",
		}
	}

	checkIfAlreadyVerified, err := NewQuery(Master(), "_User", types.M{"username": username, "emailVerified": true}, types.M{}, nil)
	if err != nil {
		return false
	}
	result, err := checkIfAlreadyVerified.Execute()
	if err != nil {
		return false
	}
	if utils.HasResults(result) {
		return true
	}

	document, err := db.Update("_User", query, updateFields, types.M{}, false)
	if err != nil {
		return false
	}
	if document == nil {
		return false
	}

	return true
}

// CheckResetTokenValidity 检查要重置密码的用户与 token 是否存在
func CheckResetTokenValidity(username, token string) types.M {
	db := orm.TalismanDBController
	// 校验 token 是否过期
	where := types.M{
		"username":          username,
		"_perishable_token": token,
	}
	if config.TConfig.PasswordPolicy && config.TConfig.ResetTokenValidityDuration > 0 {
		where["_perishable_token_expires_at"] = types.M{
			"$gt": utils.TimetoString(time.Now().UTC()),
		}
	}
	option := types.M{"limit": 1}
	results, err := db.Find("_User", where, option)
	if err != nil {
		return nil
	}
	if len(results) != 1 {
		return nil
	}

	return utils.M(results[0])
}

// UpdatePassword 更新指定用户的密码
func UpdatePassword(username, token, newPassword string) error {
	user := CheckResetTokenValidity(username, token)
	if user == nil {
		return errors.New("Invalid token")
	}

	err := updateUserPassword(user["objectId"].(string), newPassword)
	if err != nil {
		return err
	}

	// 清空重置密码 token
	db := orm.TalismanDBController
	selector := types.M{"username": username}
	update := types.M{
		"_perishable_token":            types.M{"__op": "Delete"},
		"_perishable_token_expires_at": types.M{"__op": "Delete"},
	}
	_, err = db.Update("_User", selector, update, types.M{}, false)

	return err
}

func updateUserPassword(userID, password string) error {
	_, err := Update(Master(), "_User", userID, types.M{"password": password}, nil)
	if err != nil {
		return err
	}
	return nil
}

func buildEmailLink(destination, username, token string) string {
	usernameAndToken := `token=` + token + `&username=` + username

	if config.ParseFrameURL() != "" {
		destinationWithoutHost := strings.Replace(destination, config.TConfig.ServerURL, "", -1)
		return config.ParseFrameURL() + `?link=` + url.QueryEscape(destinationWithoutHost) + `&` + usernameAndToken
	}
	return destination + `?` + usernameAndToken
}
