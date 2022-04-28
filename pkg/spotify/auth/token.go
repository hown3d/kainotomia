package auth

import (
	"time"

	"golang.org/x/oauth2"
)

func RefreshToken(current oauth2.Token) oauth2.Token {
	// current token is not expired yet, reuse it
	if current.Expiry.Before(time.Now()) {
		return current
	}
	// TODO: implement
	panic("not implemented")
}
