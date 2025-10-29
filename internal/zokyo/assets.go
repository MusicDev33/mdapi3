package zokyo

import (
	"math/rand"
)

var passwordMessages = []string{
	"Incorrect password. Please try again.",
	"Sorry, that password is incorrect. Please try again.",
	"Incorrect login credentials. Please check your password and try again.",
	"Your password is incorrect. Please try again.",
	"Sorry, we couldn't verify your password. Please try again.",
	"Your password does not match our records. Please try again.",
	"That password isn't right. Please try again.",
}

func GetPassMsg() string {
	return passwordMessages[rand.Intn(len(passwordMessages))]
}

var userNotFoundMessages = []string{
	"User not found. Please check your credentials and try again.",
	"Sorry, we couldn't find that user. Please try again.",
	"We couldn't find a user with that username or email. Please try again.",
	"The username or email you entered is not registered. Please try again.",
	"We're sorry, that user doesn't exist. Please check your input and try again.",
}

func GetUserNotFoundMsg() string {
	return userNotFoundMessages[rand.Intn(len(userNotFoundMessages))]
}
