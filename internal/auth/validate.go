package auth

import "net/mail"

const minPasswordLen = 8

func validEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}
