package models

type Token struct {
	Base
	RawToken string
	ClientID string
	ExpressIn int64
}