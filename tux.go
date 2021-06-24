package tux

import (
	"net/http"
)

type Entry interface {
	Entry()
}

type Option interface {
	Option()
}

type Group struct {
	Path     string
	Children []Entry
	Options  []Option
}

func (_ Group) Entry() {}

type Handle struct {
	Path    string
	Handler http.Handler
	Options []Option
}

func (_ Handle) Entry() {}

type Argument struct {
	Name    string
	Child   Entry
	Options []Option
}

func (_ Argument) Entry() {}

type Matcher interface {
	Match(http.ResponseWriter, *http.Request) bool
}

type MatchFunc func(http.ResponseWriter, *http.Request) bool

func (f MatchFunc) Match(w http.ResponseWriter, r *http.Request) bool {
	return f(w, r)
}

type ArgumentsHandler interface {
	http.Handler
	ServeHTTPWithArguments(http.ResponseWriter, *http.Request, map[string]string)
}

type ArgumentsHandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

func (f ArgumentsHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r, nil)
}

func (f ArgumentsHandlerFunc) ServeHTTPWithArguments(w http.ResponseWriter, r *http.Request, args map[string]string) {
	f(w, r, args)
}
