package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
)

func main() {

	mux := pat.New()
	//serve the images first
	mux.Get("/images/", http.StripPrefix("/images/",
		http.FileServer(http.Dir("images"))))
	mux.Get("/css/", http.StripPrefix("/css/",
		http.FileServer(http.Dir("css"))))
	mux.Get("/", http.HandlerFunc(home))
	mux.Post("/", http.HandlerFunc(send))
	mux.Get("/confirmation", http.HandlerFunc(confirmation))

	log.Println("Listening...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func images(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./images/")
}

func home(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/home.html", nil)
}

func send(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate form
	msg := &Message{
		Name:       r.PostFormValue("name"),
		Phone:      r.PostFormValue("phone"),
		Address:    r.PostFormValue("address"),
		BusyKidNum: r.PostFormValue("busykidnum"),
		Email:      r.PostFormValue("email"),
		Restaurant: r.PostFormValue("restaurant"),
		DateTime:   r.PostFormValue("datetime"),
		Content:    r.PostFormValue("content"),
		Errors:     map[string]string{},
	}

	if msg.Validate() == false {
		render(w, "templates/home.html", msg)
		return
	}

	// Step 2: Send contact form message in an email
	if err := msg.Deliver(); err != nil {
		log.Println(err)
		log.Println("sent email to delivery boy")
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	if err := msg.Deliver_Receipt(); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	// Step 3: Redirect to confirmation page
	//http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
	render(w, "templates/confirmation.html", msg)
}

func confirmation(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/confirmation.html", nil)
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}
