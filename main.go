package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title       string
	Student     string
	Client      string
	Description []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Description, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Description: body}, nil
}

func main() {
	p1 := &Page{Title: "TestPage", Description: []byte("This is a sample page.")}
	p1.save()
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, my name is %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/view/" {
		fmt.Fprint(w, "<h1>You didn't specify a page!")
	} else {
		title := r.URL.Path[len("/view/"):]
		p, _ := loadPage(title)
		fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Description)
	}
}
