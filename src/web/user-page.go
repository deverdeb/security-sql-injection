package web

import (
	"log"
	"net/http"
	"strconv"
	"web-appli/src/users"
)

// UserPage est la fonction de rendu de la page d'édition d'un utilisateur
func UserPage(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	Sessions.CheckAuthentication(responseWriter, httpRequest, func(user *users.User) {
		userId, err := strconv.ParseInt(httpRequest.URL.Query().Get("id"), 10, 64)
		if err != nil {
			// Echec de la récupération de l'identifiant
			displayErrorPage(responseWriter, 400, "Invalid id ["+httpRequest.URL.Query().Get("id")+"]")
		}
		if user.Id != userId && !user.IsAdmin {
			// Page non autorisée aux utilisateurs qui ne sont pas administrateur ou à l'utilisateur lui-même.
			// Retourner à la page d'index
			http.Redirect(responseWriter, httpRequest, "index", http.StatusSeeOther)
		}
		// Récupérer l'utilisateur et l'action sur l'utilisateur
		editedUser, err := users.Service.FindById(userId)
		action := httpRequest.URL.Query().Get("action")
		if err != nil {
			// Problème de base de données...
			displayErrorPageFromError(responseWriter, err)
		} else if editedUser == nil {
			// L'utilisateur demandé n'existe pas, ce qui n'est pas cohérent
			displayErrorPage(responseWriter, 404, "User is not found or forbidden")
		} else if httpRequest.Method == "POST" {
			// Mise à jour d'un utilisateur
			updateUser(editedUser, responseWriter, httpRequest, user)
		} else if action == "delete" {
			// Supprimer la tâche
			deleteUser(editedUser, responseWriter, httpRequest, user)
		} else {
			// Affichage d'un utilisateur
			displayUserPage(responseWriter, user, editedUser)
		}

	})
}

func updateUser(editedUser *users.User, responseWriter http.ResponseWriter, httpRequest *http.Request, user *users.User) {
	// Extraire les champs du formulaire et compléter la tâche
	editedUser = completeUser(editedUser, httpRequest)
	// Sauver la tâche
	err := users.Service.Save(editedUser)
	if err != nil {
		log.Printf("[ERROR] failed to update user: %v", err)
		displayErrorPageFromError(responseWriter, err)
	}
	// Rediriger vers la page d'index
	http.Redirect(
		responseWriter, httpRequest,
		"users",
		http.StatusSeeOther,
	)
}

func completeUser(editedUser *users.User, httpRequest *http.Request) *users.User {
	httpRequest.ParseForm()
	editedUser.Firstname = httpRequest.FormValue("Firstname")
	editedUser.Lastname = httpRequest.FormValue("Lastname")
	editedUser.IsAdmin = httpRequest.Form.Has("IsAdmin")
	return editedUser
}

func deleteUser(editedUser *users.User, responseWriter http.ResponseWriter, httpRequest *http.Request, user *users.User) {
	// Supprimer l'utilisateur'
	_ = users.Service.Delete(editedUser)
	// Rediriger vers la page d'index
	http.Redirect(
		responseWriter, httpRequest,
		"users",
		http.StatusSeeOther,
	)
}

// displayUserPage est la fonction de rendu de la page d'édition des utilisateurs
func displayUserPage(responseWriter http.ResponseWriter, user *users.User, editedUser *users.User) {
	err := htmlTemplates.ExecuteTemplate(responseWriter, "user.gohtml", struct {
		User       *users.User
		EditedUser *users.User
	}{
		User:       user,
		EditedUser: editedUser,
	})
	if err != nil {
		log.Printf("[ERROR] failed to build user page: %v", err)
		displayErrorPageFromError(responseWriter, err)
	}
}
