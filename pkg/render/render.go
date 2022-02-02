package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/akrambelajouza/bookings/pkg/config"
	"github.com/akrambelajouza/bookings/pkg/models"
)

var functions = template.FuncMap{}

// app points to the configuration
var app *config.AppConfig

// Newtemplates set the config for template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// RenderTempate render html files
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	/*
			_, err := RenderTemplateTest(w)
			if err != nil {
				fmt.Println("Error getting template cache: ", err)
				return
			}

		parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
		err := parsedTemplate.Execute(w, nil)
		if err != nil {
			fmt.Println("Error parsing template: ", err)
			return
		}
	*/

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	/*
		tc, err := CreateTemplateCache()

		if err != nil {
			log.Fatal((err))
		}
	*/
	t, ok := tc[tmpl]

	if !ok {
		log.Fatal(("Could not create template from template cache"))
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing to browser: ", err)
	}
}

// CreateTemplateCache creates templte cache as map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		fmt.Println("Error rendering template: ", err)
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		//		fmt.Println("Page is currently: ", page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, err
}
