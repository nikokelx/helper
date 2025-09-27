/*
Raffael Rot
V0.1.1
*/

package helper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type SessionResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
}

type WebResponse struct {
	Code    int         `json:"code,omitzero"`
	Status  string      `json:"status,omitzero"`
	Success bool        `json:"success,omitzero"`
	Data    interface{} `json:"data,omitzero"`
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func Check(values ...uuid.UUID) uuid.UUID {
	if len(values) == 0 {
		return uuid.New()
	} else {
		return values[0]
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

func WriteResponseBody(writer http.ResponseWriter, status int, response interface{}) {
	writer.WriteHeader(status)
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
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

func ParseToken(authorizationToken string, secret string) (*JWTClaims, error) {
	encodedToken := strings.TrimPrefix(authorizationToken, "Bearer ")

	token, err := jwt.ParseWithClaims(encodedToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	PanicIfError(err)

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func NewDatabaseConnection(
	username string,
	password string,
	dbname string,
	host string,
	port string,
) *sql.DB {
	// Create connection string
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		username,
		password,
		dbname,
		host,
		port,
	)

	db, errDB := sql.Open("postgres", connStr)
	PanicIfError(errDB)

	errPing := db.Ping()
	PanicIfError(errPing)

	return db
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
	} else {
		return nil
	}
}
