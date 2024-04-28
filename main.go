package main

import (
    "os"
    "log"
    "regexp"
    "html/template"
    "net/http"
    "go-wiki/src/wikifiles"
    "go-wiki/src/wikiparser"
)

var (
    templates = template.Must(template.ParseFiles("edit.html", "view.html", "dash.html"))
    validPath = regexp.MustCompile("^/(edit|save|view|dash)/([a-zA-Z0-9]+)$")
)

type Page struct {
    Title string
    Body []byte
    HTML template.HTML
}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func dashHandler(w http.ResponseWriter, r *http.Request) {
    files, err := wikifiles.WalkMatch("./", "*.txt")

    err = templates.ExecuteTemplate(w, "dash.html", files)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func createHandler(w http.ResponseWriter, r *http.Request) {
    page := r.FormValue("page")
    
    http.Redirect(w, r, "/edit/"+page, http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    body := wikiparser.ApplySyntax(p.Body)
    p.HTML = template.HTML(body)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    submit := r.FormValue("submit")
    if submit == "Cancel" {
        http.Redirect(w, r, "/view/"+title, http.StatusFound)
        return
    }
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

/**
 * Front controller - validate path
 */
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
    http.HandleFunc("/", dashHandler)
    http.HandleFunc("/create", createHandler)
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))

    log.Fatal(http.ListenAndServe(":8080", nil))
}
