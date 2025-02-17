package controllers

import "net/http"

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	if err := Tmpl.ExecuteTemplate(w, "errors.html", map[string]interface{}{
		"Status": status,
	}); err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
