package wikifiles

import (
    "os"
    "path/filepath"
)

type WikiFile struct {
    Title string
    LastUpdated string
}

func WalkMatch(root, pattern string) ([]WikiFile, error) {
    var matches []WikiFile
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
            wikiFile := WikiFile{Title: filename, LastUpdated: info.ModTime().Format("Jan 02, 2006 3:04 PM")}
            matches = append(matches, wikiFile)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return matches, nil
}
