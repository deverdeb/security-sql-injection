package main

import (
    "fmt"
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
    http.HandleFunc("/index", web.IndexPage)
    http.HandleFunc("/login", web.LoginPage)
    http.HandleFunc("/search", web.SearchPage)
    http.HandleFunc("/task", web.TaskPage)
    http.HandleFunc("/logout", web.LogoutPage)

    // Ajoutes les ressources statiques
    http.Handle("/", web.BuildHttpStaticHandler("index"))

    // Lancer le serveur
    log.Printf("Start web-appli on port %d", port)
    err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
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

var priorityList = []tasks.Priority{tasks.Urgent, tasks.Normal, tasks.Basse}
var statusList = []tasks.Status{tasks.EnAttente, tasks.AFaire, tasks.EnCours, tasks.Terminee, tasks.Abandonnee}
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func insertTasks() {
    allUsers, _ := users.Service.FindAll()
    for idx, user := range allUsers {
        nbTask := random.Intn(6) + 5
        for nbtask := 0; idx < nbTask; idx++ {
            tasks.Service.Save(&tasks.Task{
                UserId:      user.Id,
                Name:        fmt.Sprintf("Tache %d-%d", idx, nbtask),
                Description: fmt.Sprintf("Tache %d de l'utilisateur %s %s", nbtask, user.Firstname, user.Lastname),
                Priority:    priorityList[random.Intn(len(priorityList))],
                Status:      statusList[random.Intn(len(statusList))],
                Archived:    random.Intn(2) == 0,
            })
        }
    }
}
