package user

import (
	"errors"

	"github.com/sethvargo/go-password/password"
	"github.com/spf13/viper"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
)

// NewPassword generates a new password
// it does not store the password
func (u *User) NewPassword() (string, error) {
	// // u.Passwd = encrypt([]byte(pw), core.Keystore.Key())
	// pe := PwEntry{
	// 	Passwd:  encrypt([]byte(pw), core.Keystore.Key()),
	// 	Created: time.Now(),
	// }
	// u.History = append([]*PwEntry{&pe}, u.History...)
	pw, err := u.generatePassword()
	if err != nil {
		return u.Password(), err
	}
	encPw := string(encrypt([]byte(pw), core.Keystore.Key()))
	for _, h := range u.History {
		if string(h.Passwd) == encPw {
			return u.Password(), errors.New("duplicate password")
		}
	}
	return pw, nil
}

func (u *User) generatePassword() (string, error) {
	pwGen, err := password.NewGenerator(&password.GeneratorInput{
		UpperLetters: viper.GetString(cfg.PasswdRuleUpper),
		LowerLetters: viper.GetString(cfg.PasswdRuleLower),
		Digits:       viper.GetString(cfg.PasswdRuleDigit),
		Symbols:      viper.GetString(cfg.PasswdRuleSymbols),
	})
	if err != nil {
		return "", err
	}
	pwLen := viper.GetInt(cfg.PasswdRuleLength)
	numDigit := viper.GetInt(cfg.PasswdRuleNumDigits)
	numSymbols := viper.GetInt(cfg.PasswdRuleNumSymbols)
	return pwGen.Generate(pwLen, numDigit, numSymbols, false, true)
}
