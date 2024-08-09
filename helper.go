package helper

import (
    "encoding/json"
    "net/http"
    "database/sql"
)

func PanicIfError(err error) {
    if err != nil {
        panic(err)
    }
}
type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func ReadRequestBody(r *http.Request, result interface{}) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(result)
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
