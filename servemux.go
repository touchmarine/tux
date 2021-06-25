package tux

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"sort"
	"strings"
)

type ServeMuxTux struct {
	m    map[string]http.Handler
	keys []string
	pos  int
}

func ServeMux(group *Group) *ServeMuxTux {
	tux := &ServeMuxTux{
		m:   make(map[string]http.Handler),
		pos: -1,
	}

	tux.group("", group)

	keys := make([]string, 0, len(tux.m))
	for k := range tux.m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tux.keys = keys

	return tux
}

func (t *ServeMuxTux) group(parentPattern string, group *Group) {
	var pattern string

	if parentPattern == "" {
		pattern += "/"
	} else {
		pattern += parentPattern
	}

	if group.Path != "" {
		pattern += group.Path + "/"
	}

	for _, child := range group.Children {
		switch c := child.(type) {
		case *Group:
			t.group(pattern, c)
		case *Argument:
			t.argument(pattern, c)
		case *Handle:
			t.handle(pattern, c)
		case Group:
			log.Fatal("expected pointer to Group")
		case Argument:
			log.Fatal("expected pointer to Argument")
		case Handle:
			log.Fatal("expected pointer to Handle")
		default:
			panic(fmt.Sprintf("tux: unexpected entry type %T", child))
		}
	}
}

func (t *ServeMuxTux) argument(parentPattern string, argument *Argument) {
	switch c := argument.Child.(type) {
	case *Handle:
		pattern := parentPattern + c.Path
		t.m[pattern] = matchGreedy(argument.Name, parentPattern, c.Handler)
	case Handle:
		log.Fatal("expected pointer to Handle")
	default:
		panic(fmt.Sprintf("tux: unexpected entry type %T", argument.Child))
	}
}

func matchGreedy(name, parentPattern string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch h := handler.(type) {
		case ArgumentsHandler:
			args := make(map[string]string)

			p := strings.TrimPrefix(cleanPath(r.URL.Path), cleanPath(parentPattern))
			args[name] = p

			h.ServeHTTPWithArguments(w, r, args)
		default:
			handler.ServeHTTP(w, r)
		}
	})
}

// used by ServeMux, copied from net/http/server.go
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		// Fast path for common case of p being the string we want:
		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
			np = p
		} else {
			np += "/"
		}
	}
	return np
}

func (t *ServeMuxTux) handle(parentPattern string, handle *Handle) {
	pattern := parentPattern + handle.Path
	handler := handle.Handler

	if pattern == "/" {
		// catch-all pattern in servemux
		handler = notFoundCatchAll(handler)
	}

	t.m[pattern] = handler
}

// notFoundCatchAll responds with not found for all but the root path.
// The "/" pattern in ServeMux matches everything, this undoes that.
func notFoundCatchAll(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (t *ServeMuxTux) Next() bool {
	if t.pos+1 < len(t.m) {
		t.pos++
		return true
	}
	return false
}

func (t *ServeMuxTux) Pattern() string {
	return t.keys[t.pos]
}

func (t *ServeMuxTux) Handler() http.Handler {
	key := t.keys[t.pos]
	return t.m[key]
}
