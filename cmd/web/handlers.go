package main

import (
	"log"
	"net/http"
)

// Home display the home page of the sites
func (app *App) Home(w http.ResponseWriter, r *http.Request) {

	// 404 if not truly root
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}

	// Get the latest requests
	requests, err := app.DB.LatestRequests(5)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Get the latest books
	books, err := app.DB.LatestBooks(10)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Get announcements, if any
	announcement, err := app.DB.GetAnnouncement()
	if err != nil {
		log.Printf("Unable to get announcements: %s", err.Error())
	} else {
		// Display home page with books and requests + announcements
		app.RenderHTML(w, r, "home.page.html", &HTMLData{
			Requests:     requests,
			Books:        books,
			Announcement: announcement,
		})
		return
	}

	// Display home page with books and requests
	app.RenderHTML(w, r, "home.page.html", &HTMLData{
		Requests: requests,
		Books:    books,
	})
}

// About display the site information page
func (app *App) About(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "about.page.html", &HTMLData{})
}
