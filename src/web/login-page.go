package web

import (
	"fmt"
	"log"
	"net/http"
	"web-appli/src/users"
)

type loginPageData struct {
	DisplayMessage bool
	Message        string
}

// LoginPage est la fonction de rendu de la page de login
func LoginPage(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method == "POST" {
		checkAuthentication(responseWriter, httpRequest)
	} else {
		displayLoginPage(responseWriter, loginPageData{DisplayMessage: false})
	}
}

func checkAuthentication(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	pageData := loginPageData{DisplayMessage: false}
	httpRequest.ParseForm()
	login := httpRequest.FormValue("login")
	password := httpRequest.FormValue("password")
	user, err := users.Service.Authentication(login, password)
	if err != nil {
		log.Printf("ERROR - %v", err)
		pageData.Message = fmt.Sprintf("error: %s", err)
	}
	if user != nil {
		//log.Printf("login = OK -> redirect to login page")
		// Lier la session à la requête
		Sessions.Create(user, responseWriter)
		// Rediriger vers la page d'index
		http.Redirect(
			responseWriter, httpRequest,
			"index",
			http.StatusSeeOther, /*StatusMovedPermanently*/
		)
		return
	} else {
		//log.Printf("login = KO -> login page")
		pageData.DisplayMessage = true
		pageData.Message = "Invalid login or password"
	}
	displayLoginPage(responseWriter, pageData)
}

func displayLoginPage(responseWriter http.ResponseWriter, pageData loginPageData) {
	err := htmlTemplates.ExecuteTemplate(responseWriter, "login.html", pageData)
	if err != nil {
		log.Printf("[ERROR] failed to build login page: %v", err)
		displayErrorPageFromError(responseWriter, err)
	}
}
