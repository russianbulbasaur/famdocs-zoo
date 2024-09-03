package helpers

import (
	"fmt"
	"net/http"
)

func InvalidParametersError(w http.ResponseWriter, code int, err error) {
	LogToFile(err)
	http.Error(w, fmt.Sprintf("Missing paramters %d", code), http.StatusBadRequest)
}

func MissingParametersError(w http.ResponseWriter, code int, err error) {
	LogToFile(err)
	http.Error(w, fmt.Sprintf("Missing paramters %d", code), http.StatusBadRequest)
}
func FormParseError(w http.ResponseWriter, code int, err error) {
	LogToFile(err)
	http.Error(w, fmt.Sprintf("something went wrong while parsing form %d", code), http.StatusBadRequest)
}

func InternalError(w http.ResponseWriter, code int, err error) {
	LogToFile(err)
	http.Error(w, fmt.Sprintf("internal server error %d", code), http.StatusInternalServerError)
}

func WriteError(w http.ResponseWriter, code int, err error) {
	LogToFile(err)
	http.Error(w, fmt.Sprintf("write error %d", code), http.StatusInternalServerError)
}
