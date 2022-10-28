package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
	"web-appli/src/db"
	"web-appli/src/tasks"
	"web-appli/src/users"
	"web-appli/src/web"
)

// Port pour l'exécution du serveur
var port = 8080

func main() {
	log.Printf("Start web-appli")

	// Initialiser la base de données
	db.Initialize()
	defer db.Release()

	insertUsers()
	insertTasks()

	// Créer les routes vers les pages
	router := mux.NewRouter()

	// Ajoutes les ressources statiques
	router.HandleFunc("/", web.IndexPage)
	router.HandleFunc("/index", web.IndexPage)
	router.HandleFunc("/login", web.LoginPage)
	router.HandleFunc("/search", web.SearchPage)
	router.HandleFunc("/task", web.TaskPage)
	router.HandleFunc("/logout", web.LogoutPage)
	router.PathPrefix("/").Handler(web.BuildHttpStaticHandler())

	// Lancer le serveur
	log.Printf("Start web-appli on port %d", port)
	server := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("web-appli execution failed: %v", err)
	}

	log.Printf("Stop web-appli")
}

func insertUsers() {
	users.Service.Create(&users.User{Firstname: "JEAN", Lastname: "DUPONT", Login: "jdupont", IsAdmin: false}, "azerty")
	users.Service.Create(&users.User{Firstname: "MARC", Lastname: "HASSIN", Login: "mhassin", IsAdmin: false}, "hassin123")
	users.Service.Create(&users.User{Firstname: "Administrateur", Lastname: "", Login: "admin", IsAdmin: true}, "321ytreza")
	users.Service.Create(&users.User{Firstname: "MARIE", Lastname: "MARTIN", Login: "mmartin", IsAdmin: false}, "mmartin")
	users.Service.Create(&users.User{Firstname: "ANNE", Lastname: "DUPOND", Login: "adupond", IsAdmin: false}, "azerty")
}

var priorityList = []tasks.Priority{tasks.Highest, tasks.High, tasks.Medium, tasks.Low, tasks.Lowest}
var statusList = []tasks.Status{tasks.Draft, tasks.Open, tasks.InProgress, tasks.Done, tasks.Closed, tasks.Abandoned}
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func insertTasks() {
	allUsers, _ := users.Service.FindAll()
	for idxUser, user := range allUsers {
		nbTask := random.Intn(6) + 5
		for idTask := 0; idTask < nbTask; idTask++ {
			_, err := tasks.Service.Save(&tasks.Task{
				UserId:      user.Id,
				Name:        fmt.Sprintf("Tache %d-%d", idxUser, idTask),
				Description: fmt.Sprintf("Tache %d de l'utilisateur %s %s", idTask, user.Firstname, user.Lastname),
				Priority:    priorityList[random.Intn(len(priorityList))],
				Status:      statusList[random.Intn(len(statusList))],
				Archived:    random.Intn(2) == 0,
			})
			if err != nil {
				log.Fatalf("failed to create task: %v", err)
			}
		}
	}
}
