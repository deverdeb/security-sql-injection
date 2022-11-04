package web

import (
	"log"
	"net/http"
	"web-appli/src/tasks"
	"web-appli/src/users"
)

// IndexPage est la fonction de rendu de la page d'index
func IndexPage(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	Sessions.CheckAuthentication(responseWriter, httpRequest, func(user *users.User) {
		taskList, err := tasks.Service.FindByUser(user)
		if err != nil {
			displayErrorPageFromError(responseWriter, err)
		} else {
			displayIndexPage(responseWriter, user, taskList, "")
		}
	})
}

// displayIndexPage est la fonction de rendu de la page d'index
func displayIndexPage(responseWriter http.ResponseWriter, user *users.User, taskList []*tasks.Task, searchText string) {
	err := htmlTemplates.ExecuteTemplate(responseWriter, "index.html", struct {
		User       *users.User
		Tasks      []*tasks.Task
		SearchText string
	}{
		User:       user,
		Tasks:      taskList,
		SearchText: searchText,
	})
	if err != nil {
		log.Printf("[ERROR] failed to build index page: %v", err)
		displayErrorPageFromError(responseWriter, err)
	}
}
