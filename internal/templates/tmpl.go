package tmpl

import (
	"html/template"
	"net/http"
)

func NewTempl() map[string]*template.Template {

	templ := make(map[string]*template.Template)

	templ["index"] = template.Must(template.ParseFiles("internal/templates/html/index.html",
		"internal/templates/html/base.html"))

	templ["lotDescription"] = template.Must(template.ParseFiles("internal/templates/html/lotDescription.html",
		"internal/templates/html/base.html"))

	return templ
}

func RenderTemplate(w http.ResponseWriter, name string, base string, viewModel interface{}, templates map[string]*template.Template) {
	tmpl, ok := templates[name]

	if !ok {
		http.Error(w, "can't find template", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, base, viewModel)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
