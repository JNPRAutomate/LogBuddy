package logbuddy

import (
	"net/http"
	"text/template"
)

func (ws *WebServer) HomeHandler(w http.ResponseWriter, r *http.Request) {
	finalTemplate := template.New("home")

	headerAsset, _ := Asset("static/tmpl/header.tmpl")
	footerAsset, _ := Asset("static/tmpl/footer.tmpl")
	bodyAsset, _ := Asset("static/tmpl/home.tmpl")

	finalTemplate.Parse(string(headerAsset))
	finalTemplate.Parse(string(footerAsset))
	finalTemplate.Parse(string(bodyAsset))

	finalTemplate.ExecuteTemplate(w, "header", nil)
	finalTemplate.ExecuteTemplate(w, "body", nil)
	finalTemplate.ExecuteTemplate(w, "footer", nil)
}
