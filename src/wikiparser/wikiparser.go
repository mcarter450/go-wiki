package wikiparser

import "regexp"

var (
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

func ApplySyntax(body []byte) []byte {
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
