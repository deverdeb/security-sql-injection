package web

import (
    "log"
    "net/http"
)

// displayErrorPage est la fonction de rendu de la page d'erreur
func displayErrorPageFromError(responseWriter http.ResponseWriter, error error) {
    responseWriter.WriteHeader(500)
    err := htmlTemplates.ExecuteTemplate(responseWriter, "error.html", error)
    if err != nil {
        log.Printf("[ERROR] failed to build error page: %v", err)
    }
}

// displayErrorPage est la fonction de rendu de la page d'erreur
func displayErrorPage(responseWriter http.ResponseWriter, httpStatus int, errorMessage string) {
    responseWriter.WriteHeader(httpStatus)
    err := htmlTemplates.ExecuteTemplate(responseWriter, "error.html", errorMessage)
    if err != nil {
        log.Printf("[ERROR] failed to build error page: %v", err)
    }
}
