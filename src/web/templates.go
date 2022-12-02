package web

import (
	"embed"
	"html/template"
	"reflect"
)

// htmlTemplateFiles pointe vers les templates HTML de l'application.
// Nous embarquons les fichiers dans l'exécutable.
//
//go:embed html-templates
var htmlTemplateFiles embed.FS

// htmlTemplates contient l'ensemble des templates des pages de l'application
var htmlTemplates = template.Must(template.New("").
	Funcs(template.FuncMap{ // Ajouter des fonctions aux templates
		"hasField": hasField,
	}).
	ParseFS(htmlTemplateFiles, "html-templates/*"))

// hasField vérifie la présence d'un champ dans une donnée
func hasField(name string, data interface{}) bool {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	return v.FieldByName(name).IsValid()
}
