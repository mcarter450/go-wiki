package wikifiles

import (
    "os"
    "path/filepath"
)

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
