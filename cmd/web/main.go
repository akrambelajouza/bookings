package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/akrambelajouza/bookings/pkg/config"
	"github.com/akrambelajouza/bookings/pkg/handlers"
	"github.com/akrambelajouza/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const serverHost = "127.0.0.1:3000"

var app config.AppConfig

var session *scs.SessionManager

func main() {

	// change this to true in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create TemplateCache ", err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Staring the application on " + serverHost)
	// _ = http.ListenAndServe(serverHost, nil)
	srv := &http.Server{
		Addr:    serverHost,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
