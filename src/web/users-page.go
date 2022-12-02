package web

import (
	"log"
	"net/http"
	"web-appli/src/users"
)

// UsersPage est la fonction de rendu de la page des utilisateurs
func UsersPage(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	Sessions.CheckAuthentication(responseWriter, httpRequest, func(user *users.User) {
		if !user.IsAdmin {
			// Page non autorisée aux utilisateurs qui ne sont pas administrateur.
			// Retourner à la page d'index
			http.Redirect(responseWriter, httpRequest, "index", http.StatusSeeOther)
		}
		userList, err := users.Service.FindAll()
		if err != nil {
			displayErrorPageFromError(responseWriter, err)
		} else {
			displayUsersPage(responseWriter, user, userList)
		}
	})
}

// displayUsersPage est la fonction de rendu de la page des utilisateurs
func displayUsersPage(responseWriter http.ResponseWriter, user *users.User, userList []*users.User) {
	err := htmlTemplates.ExecuteTemplate(responseWriter, "users.gohtml", struct {
		User  *users.User
		Users []*users.User
	}{
		User:  user,
		Users: userList,
	})
	if err != nil {
		log.Printf("[ERROR] failed to build users page: %v", err)
		displayErrorPageFromError(responseWriter, err)
	}
}
