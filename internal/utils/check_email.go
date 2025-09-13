package utils

import "net/mail"

func CheckIsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
