package main

import (
    "os"
    "log"
    "regexp"
    "html/template"
    "net/http"
    "path/filepath"
)

var (
    templates = template.Must(template.ParseFiles("edit.html", "view.html", "dash.html"))
    validPath = regexp.MustCompile("^/(edit|save|view|dash)/([a-zA-Z0-9]+)$")
    h4Regex = regexp.MustCompile("====([^=]+)====")
    h3Regex = regexp.MustCompile("===([^=]+)===")
    h2Regex = regexp.MustCompile("==([^=]+)==")
    h1Regex = regexp.MustCompile("=([^=]+)=")
    strongRegex = regexp.MustCompile("[']{3}([^']+)[']{3}")
    emRegex = regexp.MustCompile("[']{2}([^']+)[']{2}")
    strikeRegex = regexp.MustCompile("~~([^~]+)~~")
    uRegex = regexp.MustCompile(`[+]{2}([^+]+)[+]{2}`)
    codeRegex = regexp.MustCompile(`\[code\]([^\[]+)\[/code\]`)
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

func wikiSyntax(body []byte) []byte {
    body = h4Regex.ReplaceAll(body, []byte("<h4>$1</h4>"))
    body = h3Regex.ReplaceAll(body, []byte("<h3>$1</h3>"))
    body = h2Regex.ReplaceAll(body, []byte("<h2>$1</h2>"))
    body = h1Regex.ReplaceAll(body, []byte("<h1>$1</h1>"))
    body = strongRegex.ReplaceAll(body, []byte("<strong>$1</strong>"))
    body = emRegex.ReplaceAll(body, []byte("<em>$1</em>"))
    body = strikeRegex.ReplaceAll(body, []byte("<strike>$1</strike>"))
    body = uRegex.ReplaceAll(body, []byte("<u>$1</u>"))
    body = codeRegex.ReplaceAll(body, []byte("<code>$1</code>"))
    return body
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

func WalkMatch(root, pattern string) ([]string, error) {
    var matches []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
            return err
        } else if matched {
            extension := filepath.Ext(path)
            filename := path[0:len(path)-len(extension)]
            matches = append(matches, filename)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return matches, nil
}

func dashHandler(w http.ResponseWriter, r *http.Request) {
    files, err := WalkMatch("./", "*.txt")

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
    body := wikiSyntax(p.Body)
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
