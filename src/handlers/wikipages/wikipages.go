package wikipages

import (
    "os"
    "strings"
    "html/template"
    "net/http"
    "go-wiki/src/wikifiles"
    "go-wiki/src/wikiparser"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "dash.html"))

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

func DashHandler(w http.ResponseWriter, r *http.Request) {
    files, err := wikifiles.WalkMatch("./", "*.txt")

    err = templates.ExecuteTemplate(w, "dash.html", files)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
    page := r.FormValue("page")
    page = strings.ReplaceAll(page, " ", "") // Strip whitespace chars
    
    http.Redirect(w, r, "/edit/"+page, http.StatusFound)
}

func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    body := wikiparser.ApplySyntax(p.Body)
    p.HTML = template.HTML(body)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    if body == "" {
        http.Error(w, "Empty page not allowed", http.StatusInternalServerError)
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
