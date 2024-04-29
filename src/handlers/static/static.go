package static

import (
    "regexp"
    "net/http"
)

var validStaticFile = regexp.MustCompile(".(css|js|jpe?g|gif|a?png|svg|webp)$")

/**
 * Detect mime type and serve file
 */
func StaticHandler(fs http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validStaticFile.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }

        contentType := "text/css"

        switch m[1] {
            case "js":
                contentType = "text/javascript"
            case "apng":
                contentType = "image/apng"
            case "gif":
                contentType = "image/gif"
            case "jpg":
                contentType = "image/jpeg"
            case "jpeg":
                contentType = "image/jpeg"
            case "png":
                contentType = "image/png"
            case "svg":
                contentType = "image/svg+xml"
            case "webp":
                contentType = "image/webp"
            default:
                contentType = "text/css"
        }

        w.Header().Add("Content-Type", contentType)
        fs.ServeHTTP(w, r)
    }
}
