package utils

import "github.com/brianvoe/gofakeit/v6"

func RandomAppName() string {
	return gofakeit.AppName()
}
