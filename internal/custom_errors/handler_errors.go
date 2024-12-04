package custom_errors

import (
	"fmt"
	"net/http"
)

func HandleServerError(w http.ResponseWriter, message string, err error) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "%v: %v", message, err)
}
