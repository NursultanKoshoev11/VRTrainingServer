package httpserver

import "net/http"

type apiError struct {
    Error apiErrorBody `json:"error"`
}

type apiErrorBody struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
    writeJSON(w, status, apiError{Error: apiErrorBody{Code: code, Message: message}})
}
