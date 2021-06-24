package main

import (
	"github.com/touchmarine/tux"
	"log"
	"net/http"
)

func main() {
	group := &tux.Group{"", []tux.Entry{
		&tux.Handle{"a", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("a"))
		}), nil},

		&tux.Group{"b", []tux.Entry{
			&tux.Group{"a", []tux.Entry{
				&tux.Handle{"a", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("b/a/a"))
				}), nil},
			}, nil},

			&tux.Handle{"b", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("b/b"))
			}), nil},
		}, nil},

		&tux.Group{"c", []tux.Entry{
			&tux.Argument{
				"id",
				&tux.Handle{"", tux.ArgumentsHandlerFunc(func(w http.ResponseWriter, r *http.Request, args map[string]string) {
					w.Write([]byte("c/" + args["id"]))
				}), nil},
				nil,
			},
		}, nil},
	}, nil}

	muxTux := tux.ServeMux(group)

	mux := http.NewServeMux()

	for muxTux.Next() {
		pattern := muxTux.Pattern()
		handler := muxTux.Handler()

		log.Printf("register; pattern=%s", pattern)

		mux.Handle(pattern, handler)
	}

	log.Fatal(http.ListenAndServe(":8082", mux))
}
