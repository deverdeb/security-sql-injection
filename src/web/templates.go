package web

import (
    "embed"
    "html/template"
)

// htmlTemplateFiles pointe vers les templates HTML de l'application.
// Nous embarquons les fichiers dans l'exécutable.
//go:embed html-templates
var htmlTemplateFiles embed.FS

// htmlTemplates contient l'ensemble des templates des pages de l'application
var htmlTemplates = template.Must(template.ParseFS(htmlTemplateFiles, "html-templates/*"))
