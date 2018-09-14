package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

//Page : every ticket created will follow this structure.
type Page struct {
	Title   string
	Body    string
	Student string
	Client  string
}

//Data : used for the viewIndex function.
type Data struct {
	Items []string
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9-_]+)$")

func (p *Page) save() error {
	filename := p.Title + ".json"
	var page Page
	json.Unmarshal([]byte(filename), &page)
	m := Page{p.Title, p.Body, p.Student, p.Client}
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}

	return ioutil.WriteFile("files/"+filename, b, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "files/" + title + ".json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var page Page
	json.Unmarshal([]byte(body), &page)
	return &Page{Title: title, Body: page.Body, Student: page.Student, Client: page.Client}, nil
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	fmt.Println("Starting Listener 🤯")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, my name is %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/view/" {
		viewIndex(w, r)
	} else {
		title := r.URL.Path[len("/view/"):]
		p, err := loadPage(title)
		if err != nil {
			http.Redirect(w, r, "/edit/"+title, http.StatusFound)
			return
		}
		t, _ := template.ParseFiles("templates/view.html")
		p.Body = template.HTMLEscapeString(p.Body)
		t.Execute(w, p)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	title = r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("templates/edit.html")
	t.Execute(w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	title = r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	student := r.FormValue("student")
	client := r.FormValue("client")
	p := &Page{Title: title, Body: body, Student: student, Client: client}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}

func viewIndex(w http.ResponseWriter, r *http.Request) {
	var data []string
	files, err := ioutil.ReadDir("./files/")
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		result := strings.TrimSuffix(f.Name(), ".json")
		data = append(data, result)
	}
	t, _ := template.ParseFiles("templates/viewIndex.html")
	p := &Data{Items: data}
	t.Execute(w, p)
}
