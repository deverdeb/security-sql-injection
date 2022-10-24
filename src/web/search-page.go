package web

import (
    "net/http"
    "web-appli/src/tasks"
    "web-appli/src/users"
)

type searchPageData struct {
    DisplayMessage bool
    Message        string
}

// LoginPage est la fonction de rendu de la page de login
func SearchPage(responseWriter http.ResponseWriter, httpRequest *http.Request) {
    Sessions.CheckAuthentication(responseWriter, httpRequest, func(user *users.User) {
        searchText := extractSearchText(httpRequest)
        taskList, err := tasks.Service.SearchByText(user, searchText)
        if err != nil {
            displayErrorPageFromError(responseWriter, err)
        } else {
            displayIndexPage(responseWriter, user, taskList, searchText)
        }
    })
}

func extractSearchText(httpRequest *http.Request) string {
    httpRequest.ParseForm()
    return httpRequest.FormValue("search")
}
