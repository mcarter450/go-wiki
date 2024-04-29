package main

import (
    "log"
    "regexp"
    "net/http"
    "go-wiki/src/handlers/static"
    "go-wiki/src/handlers/wikipages"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([^/]+)$")

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
    http.HandleFunc("/", wikipages.DashHandler)
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", static.StaticHandler(fs)))
    http.HandleFunc("/create", wikipages.CreateHandler)
    http.HandleFunc("/view/", makeHandler(wikipages.ViewHandler))
    http.HandleFunc("/edit/", makeHandler(wikipages.EditHandler))
    http.HandleFunc("/save/", makeHandler(wikipages.SaveHandler))

    log.Fatal(http.ListenAndServe(":8080", nil))
}
