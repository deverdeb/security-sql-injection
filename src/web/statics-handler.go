package web

import (
    "embed"
    "io/fs"
    "log"
    "net/http"
    "os"
    "strings"
)

// staticsEmbeddedFiles pointe vers les fichiers statiques de l'application.
// Nous embarquons les fichiers dans l'exécutable.
//go:embed public
var staticEmbeddedFiles embed.FS

// httpStaticsHandler est le handler HTTP permettant de récupérer les fichiers statiques.
type httpStaticsHandler struct {
    // staticFilesHandler est le handler HTTP permettant de récupérer les fichiers statiques.
    staticFilesHandler http.Handler
    // Homepage pour les redirections
    home string
}

func BuildHttpStaticHandler(home string) http.Handler {
    // Ajoutes les ressources statiques
    staticFilesFs, err := fs.Sub(fs.FS(staticEmbeddedFiles), "public")
    if err != nil {
        log.Fatalf("initialize static files failed: %v", err)
    }
    fileSystem := http.FS(staticFilesFs)
    protectedFileSystem := neuteredFileSystem{fileSystem}
    staticFilesHandler := http.FileServer(protectedFileSystem)

    return httpStaticsHandler{
        staticFilesHandler: staticFilesHandler,
        home:               home,
    }
}

func (httpHandler httpStaticsHandler) ServeHTTP(responseWriter http.ResponseWriter, httpRequest *http.Request) {
    //log.Printf("URI: %s", httpRequest.RequestURI)
    //log.Printf("URL: %s", httpRequest.URL)
    //log.Printf("Host: %s", httpRequest.Host)
    //log.Printf("RemoteAddr: %s", httpRequest.RemoteAddr)

    uri := strings.TrimSpace(httpRequest.RequestURI)
    if uri == "" || uri == "/" {
        // Redirection vers la home page
        http.Redirect(
            responseWriter, httpRequest,
            httpHandler.home,
            http.StatusMovedPermanently,
        )
    }

    httpHandler.staticFilesHandler.ServeHTTP(responseWriter, httpRequest)
}

// Par défaut, le http.FileSystem permet l'accès à un niveau dossier et remonte alors la liste des fichiers.
// neuteredFileSystem est utilisé pour éviter que le FileSystem ne remonte la liste des fichiers d'un dossier.
type neuteredFileSystem struct {
    // httpFileSystem est le http.FileSystem dans lequel récupérer les fichiers statiques.
    httpFileSystem http.FileSystem
}

// Open fait la même chose que http.FileSystem.Open, mais en ne remontant pas les dossiers.
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
    // Vérifier si c'est un dossier et bloqué dans ce cas
    file, err := nfs.httpFileSystem.Open(path)
    if err != nil {
        return nil, err
    }
    infoFile, err := file.Stat()
    if err != nil {
        return nil, err
    }
    if infoFile.IsDir() {
        // C'est un dossier, bloquer (en remontant une erreur)
        return nil, os.ErrNotExist
    }

    // Ce n'est pas un dossier
    return file, nil
}
