package repository

import (
	"errors"
	"regexp"
)

// Validasi nomor telepon
func ValidatePhoneNumber(phoneNumber string) error {
	regex := `^(62|0)8[1-9][0-9]{6,9}$`
	match, _ := regexp.MatchString(regex, phoneNumber)
	if !match {
		return errors.New("nomor telepon tidak valid")
	}
	return nil
}

// Ubah jika user menginput nomor telepon dengan 08 menjadi +628
func ChangePhoneNumber(phoneNumber string) string {
	regex := `^08[1-9][0-9]{6,9}$`
	match, _ := regexp.MatchString(regex, phoneNumber)
	if match {
		phoneNumber = "62" + phoneNumber[1:]
	}
	return phoneNumber
}

// Validasi email
func ValidateEmail(email string) error {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	match, _ := regexp.MatchString(regex, email)
	if !match {
		return errors.New("email tidak valid")
	}
	return nil
}