package web

import (
	"github.com/google/uuid"
	"net/http"
	"time"
	"web-appli/src/users"
)

const CookieName = "session-id"
const MaxSessionDurationInSecond = 300. // Session de 5 minutes

// SessionsMap est une map contenant les sessions sur le serveur
type SessionsMap struct {
	sessionsById map[string]*Session
	loginPage    string
}

// Session sur le serveur
type Session struct {
	id         string
	lastAccess time.Time
	user       *users.User
	loginPage  string
}

// Sessions du serveur
var Sessions = SessionsMap{
	sessionsById: make(map[string]*Session),
	loginPage:    "login",
}

// DurationFromLastAccess indique la durée depuis le dernier accès depuis la session
func (session Session) DurationFromLastAccess() time.Duration {
	return time.Since(session.lastAccess)
}

// IsExpired indique si la session est expirée
func (session Session) IsExpired() bool {
	return session.DurationFromLastAccess().Seconds() > MaxSessionDurationInSecond
}

// IsExpired indique si la session est expirée
func (session Session) ExpirationTime() time.Time {
	return session.lastAccess.Add(MaxSessionDurationInSecond * time.Second)
}

func (session Session) LinkToResponse(responseWriter http.ResponseWriter) {
	sessionCookie := http.Cookie{
		Name:    CookieName,
		Value:   session.id,
		Expires: session.ExpirationTime(),
	}
	http.SetCookie(responseWriter, &sessionCookie)
}

// Get permet de récupérer une session.
// Cela rafraichit le dernier accès à la session.
func (sessions SessionsMap) Get(id string) *Session {
	session := sessions.sessionsById[id]
	if session != nil {
		// des sessions expirées : nettoyage
		if session.IsExpired() {
			sessions.Clean()
		} else {
			// Mettre à jour le dernier accès
			session.lastAccess = time.Now()
			// Retourner l'utilisateur
			return session
		}
	}
	return nil
}

// Delete permet de supprimer une session.
func (sessions SessionsMap) Delete(id string) {
	session := sessions.sessionsById[id]
	if session != nil {
		delete(sessions.sessionsById, session.id)
	}
}

// Clean nettoie les sessions expirées
func (sessions SessionsMap) Clean() {
	// Récupérer les sessions obsolètes
	toRemoved := make([]string, 0, 5)
	for id, session := range sessions.sessionsById {
		if session.IsExpired() {
			toRemoved = append(toRemoved, id)
		}
	}
	// Virer les sessions
	for _, id := range toRemoved {
		delete(sessions.sessionsById, id)
	}
}

// Create permet de créer une session pour un utilisateur et l'associe à la réponse
func (sessions SessionsMap) Create(user *users.User, responseWriter http.ResponseWriter) *Session {
	session := &Session{
		id:         uuid.New().String(),
		lastAccess: time.Now(),
		user:       user,
	}
	sessions.sessionsById[session.id] = session
	session.LinkToResponse(responseWriter)
	return session
}

// Check vérifie la session et retourne l'utilisateur lié
func (sessions SessionsMap) Check(responseWriter http.ResponseWriter, httpRequest *http.Request) *users.User {
	cookie, err := httpRequest.Cookie(CookieName)
	if err != nil || cookie == nil {
		// Pas de cookies
		return nil
	}
	id := cookie.Value
	session := sessions.Get(id)
	if session == nil {
		// Pas de session
		return nil
	}
	// Relier la session à la réponse
	session.LinkToResponse(responseWriter)
	// Retourner l'utilisateur
	return session.user
}

// CheckAuthentication vérifie la session et l'authentification.
// Si la session est trouvée avec un utilisateur, exécuter la fonction en paramètre.
// Si la session n'est pas trouvée, rediriger vers la page de login.
func (sessions SessionsMap) CheckAuthentication(responseWriter http.ResponseWriter, httpRequest *http.Request, funcIfOk func(*users.User)) {
	user := sessions.Check(responseWriter, httpRequest)
	if user == nil {
		// Utilisateur non connecté
		// Rediriger vers la page de login
		http.Redirect(
			responseWriter, httpRequest,
			sessions.loginPage,
			http.StatusSeeOther,
		)
	} else {
		funcIfOk(user)
	}
}

// Logout supprime la session de l'utilisateur
func (sessions SessionsMap) Logout(httpRequest *http.Request) {
	cookie, err := httpRequest.Cookie(CookieName)
	if err != nil || cookie == nil {
		// Pas de cookies - pas de session - rien à faire
		return
	}
	id := cookie.Value
	sessions.Delete(id)
}
