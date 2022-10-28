package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
)

// staticsEmbeddedFiles pointe vers les fichiers statiques de l'application.
// Nous embarquons les fichiers dans l'exécutable.
//
//go:embed public
var staticEmbeddedFiles embed.FS

// BuildHttpStaticHandler met en place le Handler HTTP remontant les ressources statiques
func BuildHttpStaticHandler() http.Handler {
	// Ajoutes les ressources statiques
	staticFilesFs, err := fs.Sub(fs.FS(staticEmbeddedFiles), "public")
	if err != nil {
		log.Fatalf("initialize static files failed: %v", err)
	}
	fileSystem := http.FS(staticFilesFs)
	protectedFileSystem := neuteredFileSystem{fileSystem}
	return http.FileServer(protectedFileSystem)
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
