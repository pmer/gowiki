package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Todo Title")
	}
	return m[2], nil
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, txt *Todo) {
	err := templates.ExecuteTemplate(w, tmpl+".html", txt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	txt, err := StoreGetTodo(title)
	if err != nil {
		log.Fatal(err)
		return
	}
	if txt == nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", txt)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	txt, err := StoreGetTodo(title)
	if err != nil {
		log.Fatal(err)
	}
	if txt == nil {
		txt = &Todo{Title: title}
	}
	renderTemplate(w, "edit", txt)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	notes := r.FormValue("notes")
	txt := Todo{Title: title, Notes: notes}
	err := StoreSetTodo(txt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	fmt.Println("Starting :)")

	err := StoreConstruct()
	if err != nil {
		log.Fatal(err)
	}

	defer StoreDestroy()

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
