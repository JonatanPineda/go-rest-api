package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func decodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func encodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func respond(w http.ResponseWriter, r *http.Request,
	status int, data interface{},
) {
	w.WriteHeader(status)
	if data != nil {
		encodeBody(w, r, data)
	}
}
func respondErr(w http.ResponseWriter, r *http.Request,
	status int, args ...interface{},
) {
	respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}
func respondHTTPErr(w http.ResponseWriter, r *http.Request,
	status int,
) {
	respondErr(w, r, status, http.StatusText(status))
}

// PathSeparator is the character used to separate
// HTTP paths.
const PathSeparator = "/"

// Path represents the path of a request.
type Path struct {
	Path string
	ID   string
}

// NewPath makes a new Path from the specified
// path string.
func NewPath(p string) *Path {
	var id string
	p = strings.Trim(p, PathSeparator)
	s := strings.Split(p, PathSeparator)
	if len(s) > 1 {
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], PathSeparator)
	}
	return &Path{Path: p, ID: id}
}

// HasID gets whether this path has an ID
// or not.
func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
