package ports

import "net/http"

type Presenter interface {
	JSON(w http.ResponseWriter, r *http.Request, v interface{})
	Error(w http.ResponseWriter, r *http.Request, err error)
}
