package helper

import (
	"database/sql"
	"encoding/json"
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

func VerifySession(url string) string {
	response, err := http.Get(url)
	PanicIfError(err)

	session := new(SessionResponse)

	ReadResponseBody(response, &session)

	return session.Data.Token
}
