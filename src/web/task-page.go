package web

import (
	"log"
	"net/http"
	"strings"
	"web-appli/src/tasks"
	"web-appli/src/users"
)

// TaskPage est la fonction de rendu de la page d'édition des tâches
func TaskPage(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	Sessions.CheckAuthentication(responseWriter, httpRequest, func(user *users.User) {
		taskId := httpRequest.URL.Query().Get("id")
		if strings.TrimSpace(taskId) == "" || strings.TrimSpace(taskId) == "0" {
			// Pas d'identifiant, nous sommes sur une création
			if httpRequest.Method == "POST" {
				// Nous sommes sur une sauvegarde, créer la tâche
				createTask(responseWriter, httpRequest, user)
			} else {
				// Nous sommes sur l'affichage du formulaire, afficher un formulaire vide
				displayTaskPage(responseWriter, user, &tasks.Task{
					UserId:      user.Id,
					Name:        "",
					Description: "",
					Priority:    tasks.Medium,
					Status:      tasks.Draft,
					Archived:    false,
				})
			}
			return
		}
		// Récupérer la tâche et l'action sur la tâche
		task, err := tasks.Service.FindByIdAndUser(taskId, user)
		action := httpRequest.URL.Query().Get("action")
		if err != nil {
			// Problème de base de données...
			displayErrorPageFromError(responseWriter, err)
		} else if task == nil {
			// La tâche demandée n'existe pas, ce qui n'est pas cohérent
			displayErrorPage(responseWriter, 404, "Resource is not found or forbidden")
		} else if httpRequest.Method == "POST" {
			// Mise à jour d'une tâche
			updateTask(task, responseWriter, httpRequest, user)
		} else if action == "delete" {
			// Supprimer la tâche
			deleteTask(task, responseWriter, httpRequest, user)
		} else {
			// Affichage d'une tâche
			// Nous sommes sur l'affichage du formulaire, afficher un formulaire vide
			displayTaskPage(responseWriter, user, task)
		}

	})
}

func createTask(responseWriter http.ResponseWriter, httpRequest *http.Request, user *users.User) {
	// Créer la nouvelle tâche
	task := &tasks.Task{
		UserId: user.Id,
	}
	// Extraire les champs du formulaire et compléter la tâche
	task = completeTask(task, httpRequest)
	// Sauver la tâche
	task, err := tasks.Service.Save(task)
	if err != nil {
		log.Printf("[ERROR] failed to create task page: %v", err)
		displayErrorPageFromError(responseWriter, err)
	}
	// Rediriger vers la page d'index
	http.Redirect(
		responseWriter, httpRequest,
		"index",
		http.StatusSeeOther,
	)
}

func updateTask(task *tasks.Task, responseWriter http.ResponseWriter, httpRequest *http.Request, user *users.User) {
	// Vérifier si nous avons le droit de modifier la tâche
	if task.UserId != user.Id {
		// La tâche demandée n'appartient pas à l'utilisateur
		displayErrorPage(responseWriter, 403, "La ressource demandée n'existe pas ou n'est pas accessible")
	}
	// Extraire les champs du formulaire et compléter la tâche
	task = completeTask(task, httpRequest)
	// Sauver la tâche
	task, err := tasks.Service.Save(task)
	if err != nil {
		log.Printf("[ERROR] failed to update task page: %v", err)
		displayErrorPageFromError(responseWriter, err)
	}
	// Rediriger vers la page d'index
	http.Redirect(
		responseWriter, httpRequest,
		"index",
		http.StatusSeeOther,
	)
}

func deleteTask(task *tasks.Task, responseWriter http.ResponseWriter, httpRequest *http.Request, user *users.User) {
	// Vérifier si nous avons le droit de supprimer la tâche
	if task.UserId != user.Id {
		// La tâche demandée n'appartient pas à l'utilisateur
		displayErrorPage(responseWriter, 403, "Resource is not found or forbidden")
	}
	// Supprimer la tâche
	_ = tasks.Service.Delete(task)
	// Rediriger vers la page d'index
	http.Redirect(
		responseWriter, httpRequest,
		"index",
		http.StatusSeeOther,
	)
}

func completeTask(task *tasks.Task, httpRequest *http.Request) *tasks.Task {
	httpRequest.ParseForm()
	task.Name = httpRequest.FormValue("Name")
	task.Description = httpRequest.FormValue("Description")
	task.Status = tasks.Status(httpRequest.FormValue("Status"))
	task.Priority = tasks.Priority(httpRequest.FormValue("Priority"))
	task.Archived = httpRequest.Form.Has("Archived")
	return task
}

// displayTaskPage est la fonction de rendu de la page d'édition des tâches
func displayTaskPage(responseWriter http.ResponseWriter, user *users.User, task *tasks.Task) {
	err := htmlTemplates.ExecuteTemplate(responseWriter, "task.html", struct {
		User *users.User
		Task *tasks.Task
	}{
		User: user,
		Task: task,
	})
	if err != nil {
		log.Printf("[ERROR] failed to build task page: %v", err)
		displayErrorPageFromError(responseWriter, err)
	}
}
