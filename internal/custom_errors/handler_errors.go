package custom_errors

import (
	"fmt"
	"net/http"
)

func ReturnErrorWithMessage(w http.ResponseWriter, message string, err error, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%v: %v", message, err)
}
