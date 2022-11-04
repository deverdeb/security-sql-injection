package web

import "github.com/gorilla/mux"

// InitializeRouter initialise les différentes routes vers les pages de l'application.
func InitializeRouter(router *mux.Router) {
	// Définir les chemins vers les pages
	router.HandleFunc("/", IndexPage)
	router.HandleFunc("/index", IndexPage)
	router.HandleFunc("/login", LoginPage)
	router.HandleFunc("/search", SearchPage)
	router.HandleFunc("/task", TaskPage)
	router.HandleFunc("/users", UsersPage)
	router.HandleFunc("/user", UserPage)
	router.HandleFunc("/logout", LogoutPage)

	// Définir les ressources statiques
	router.PathPrefix("/").Handler(BuildHttpStaticHandler())
}
