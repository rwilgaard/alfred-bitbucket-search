package util

import (
	"regexp"
	"strings"

	aw "github.com/deanishe/awgo"
)

type ParsedQuery struct {
    Text     string
    Projects []string
}

var (
    BackIcon = &aw.Icon{Value: "icons/go-back.png"}
)

func ParseQuery(query string) *ParsedQuery {
    q := new(ParsedQuery)
    projectRegex := regexp.MustCompile(`^@\w+`)

    for _, w := range strings.Split(query, " ") {
        switch {
        case projectRegex.MatchString(w):
            q.Projects = append(q.Projects, w[1:])
        default:
            q.Text = q.Text + w + " "
        }
    }

    return q
}

func CheckQueryForAutoCompletion(query string) string {
    for _, w := range strings.Split(query, " ") {
        switch w {
        case "@":
            return "projects"
        }
    }
    return ""
}
