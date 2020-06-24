package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/forms"
	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/models"
	"github.com/gorilla/mux"
)

// Home page of site
func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	// 404 if not truly root
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}

	// Get the latest requests
	requests, err := app.Request.LatestRequests()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "home.page.html", &HTMLData{
		Requests: requests,
	})
}

// ShowRequest displays a single request
func (app *App) ShowRequest(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Get requested snippet id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	// Get request
	request, err := app.Request.GetRequest(id)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if request == nil {
		app.NotFound(w)
		return
	}

	// Get the previous flashes, if any.
	if flashes := session.Flashes("default"); len(flashes) > 0 {
		// Save session
		err = session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.RenderHTML(w, r, "showrequest.page.html", &HTMLData{
			Request: request,
			Flash:   fmt.Sprintf("%v", flashes[0]),
		})
	} else {
		app.RenderHTML(w, r, "showrequest.page.html", &HTMLData{
			Request: request,
			Flash:   "",
		})
	}
}

// NewRequest displays the new request form
func (app *App) NewRequest(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "newrequest.page.html", &HTMLData{
		Form: &forms.NewRequest{},
	})
}

// CreateRequest creates a new request
func (app *App) CreateRequest(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Get requester
	valid, user := app.LoggedIn(r)
	if !valid {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Model the new request based on html form
	form := &forms.NewRequest{
		Requester: user.Username,
		Title:     r.PostForm.Get("title"),
	}

	// Validate form
	if !form.Valid() {
		app.RenderHTML(w, r, "newrequest.page.html", &HTMLData{Form: form})
		return
	}

	// Insert the new request
	id, err := app.Request.InsertRequest(form.Requester, form.Title, r.RemoteAddr)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Save success message
	session.AddFlash("Your request was saved successfully!", "default")

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/request/%d", id), http.StatusSeeOther)
}

// SignupUser presents a form to gather user information
func (app *App) SignupUser(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "signup.page.html", &HTMLData{
		Form: &forms.NewUser{},
	})
}

// CreateUser uses a form to create a user account
func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Model the new user based on html form
	form := &forms.NewUser{
		Username: r.PostForm.Get("username"),
		Email:     r.PostForm.Get("email"),
		Password:     r.PostForm.Get("password"),
	}

	// Validate form
	if !form.Valid() {
		app.RenderHTML(w, r, "signup.page.html", &HTMLData{Form: form})
		return
	}

	// Insert the new user
	err = app.User.InsertUser(form.Username, form.Email, form.Password)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Save success message
	session.AddFlash("Your account was created successfully! Please login.", "default")

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// LoginUser 
func (app *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Get the previous flashes, if any.
	if flashes := session.Flashes("default"); len(flashes) > 0 {
		// Save session
		err := session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.RenderHTML(w, r, "login.page.html", &HTMLData{
			Form: &forms.NewUser{},
			Flash:   fmt.Sprintf("%v", flashes[0]),
		})
	} else {
		app.RenderHTML(w, r, "login.page.html", &HTMLData{
			Form: &forms.NewUser{},
			Flash:   "",
		})
	}
}

// VerifyUser
func (app *App) VerifyUser(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	user := &models.User{}
	user, err = app.User.AuthenticateUser(r.PostForm.Get("username"), r.PostForm.Get("password"))

	if user == (&models.User{}) {
		// Save failure message
		session.AddFlash("Invalid Login", "default")

		// Save session
		err = session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}

	// Save user info
	session.Values["user"] = user

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/request/new", http.StatusSeeOther)
}

// LogoutUser
func (app *App) LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	session.Values["user"] = &models.User{}

	// Save session
	err := session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}