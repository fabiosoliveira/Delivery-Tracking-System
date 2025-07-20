package utils

import (
	"errors"
	"html/template"
	"net/http"
	"strings"
)

func TrowError(err error, w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("template/error.gohtml"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	tpl.Execute(w, struct {
		Error string
		Url   string
	}{err.Error(), r.URL.String()})
}

func ErrorsJoin(errs ...error) error {

	var nonNilErrors []string
	for _, err := range errs {
		if err != nil {
			nonNilErrors = append(nonNilErrors, err.Error())
		}
	}
	if len(nonNilErrors) == 0 {
		return nil
	}
	return errors.New(strings.Join(nonNilErrors, ", "))

}
