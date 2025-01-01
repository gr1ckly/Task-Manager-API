package server

import (
	"net/http"
	"strconv"
)

func buildErrorJSON(err error) []byte {
	errMessage := "{\"error\":\"" + err.Error() + "\"}"
	return []byte(errMessage)
}

func buildResultJSON(body []byte) []byte {
	resultMessage := "{\"result\":\"" + string(body) + "\"}"
	return []byte(resultMessage)
}

func sendError(respWriter *http.ResponseWriter, err error, status int) error {
	(*respWriter).WriteHeader(status)
	errMessage := buildErrorJSON(err)
	(*respWriter).Header().Set("Content-Type", "application/json")
	(*respWriter).Header().Set("Content-Length", strconv.Itoa(len(errMessage)))
	_, writeErr := (*respWriter).Write([]byte(errMessage))
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func sendResult(respWriter *http.ResponseWriter, body []byte) error {
	(*respWriter).WriteHeader(200)
	resultMessage := buildResultJSON(body)
	(*respWriter).Header().Set("Content-Type", "application/json")
	(*respWriter).Header().Set("Content-Length", strconv.Itoa(len(resultMessage)))
	_, writeErr := (*respWriter).Write([]byte(resultMessage))
	if writeErr != nil {
		return writeErr
	}
	return nil
}
