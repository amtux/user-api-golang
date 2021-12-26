package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/pkg/errors"
)

/* Resources:
https://pkg.go.dev/github.com/go-chi/jwtauth@v1.2.0#section-documentation
https://github.com/titpetric/books/tree/master/api-foundations/chapter4b-jwt */

var jwtExpiry = time.Hour * 24

type JWT struct {
	tokenClaim string
	tokenAuth  *jwtauth.JWTAuth
}

func (JWT) New(secret string) *JWT {
	return &JWT{
		tokenClaim: "user_id",
		tokenAuth:  jwtauth.New("HS256", []byte(secret), nil),
	}
}

// abstract out verifier
func (jwt *JWT) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(jwt.tokenAuth)
}

// abstract out authenticator
func (jwt *JWT) Authenticator(h http.Handler) http.Handler {
	return jwtauth.Authenticator(h)
}

func (jwt *JWT) Decode(r *http.Request) (int, error) {
	token, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || token == nil {
		return 0, errors.Wrap(err, "Empty or invalid JWT")
	}
	return int(claims[jwt.tokenClaim].(float64)), nil
}

func (jwt *JWT) Encode(id int) (string, error) {
	claims := map[string]interface{}{jwt.tokenClaim: id}
	jwtauth.SetExpiry(claims, time.Now().Add(jwtExpiry))
	jwtauth.SetIssuedNow(claims)
	_, tokenString, err := jwt.tokenAuth.Encode(claims)
	if err != nil {
		msg := fmt.Sprintf("Failed to encode jwt token for user: %v", id)
		return "", errors.Wrap(err, msg)
	}
	return tokenString, nil
}
