package tux_test

import (
	"github.com/touchmarine/tux"
	"log"
	"net/http"
)

func Example() {
	group := &tux.Group{"", []tux.Entry{
		&tux.Handle{"", http.HandlerFunc(handleAll), nil},
		&tux.Handle{"a", http.HandlerFunc(handleA), nil},
	}, nil}

	muxTux := tux.ServeMux(group)

	mux := http.NewServeMux()

	for muxTux.Next() {
		pattern := muxTux.Pattern()
		handler := muxTux.Handler()

		mux.Handle(pattern, handler)
	}

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("all"))
}

func handleA(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("a"))
}

func allowMethods(methods ...string) tux.Matcher {
	return tux.MatchFunc(func(w http.ResponseWriter, r *http.Request) bool {
		for _, method := range methods {
			if method == r.Method {
				return true
			}
		}
		return false
	})
}
