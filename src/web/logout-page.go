package web

import (
    "net/http"
)

// LogoutPage est la fonction de rendu de la page de déconnexion
func LogoutPage(responseWriter http.ResponseWriter, httpRequest *http.Request) {
    Sessions.Logout(httpRequest)
    // Rediriger vers la page de login
    http.Redirect(
        responseWriter, httpRequest,
        "login",
        http.StatusSeeOther,
    )
}
