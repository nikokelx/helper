/*
Raffael Rot
V0.1.1
*/

package helper

import (
	"database/sql"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
)

type SessionResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
}

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadRequestBody(r *http.Request, result interface{}) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(result)
	PanicIfError(err)
}

func ReadResponseBody(r *http.Response, result interface{}) {
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(result)
	PanicIfError(err)
}

func WriteResponseBody(write http.ResponseWriter, response interface{}) {
	write.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(write)
	err := encoder.Encode(response)
	PanicIfError(err)
}

func CommitOrRollback(tx *sql.Tx) {
	err := recover()

	if err != nil {
		errRollback := tx.Rollback()
		PanicIfError(errRollback)
		panic(err)
	} else {
		errCommit := tx.Commit()
		PanicIfError(errCommit)
	}
}

type JWTClaims struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func VerifySession(url string) *JWTClaims {
	response, err := http.Get(url)
	PanicIfError(err)

	session := new(SessionResponse)

	ReadResponseBody(response, &session)

	token, err := jwt.ParseWithClaims(session.Data.Token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("my-secret-key"), nil
	})
	PanicIfError(err)
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims
	}
}
