package tux_test

import (
	"fmt"
	"github.com/touchmarine/tux"
	"net/http"
	"testing"
)

func TestServeMux(t *testing.T) {
	cases := []struct {
		name string
		say  string
	}{
		{
			"",
			"all",
		},
		{
			"a",
			"a",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			group := &tux.Group{"", []tux.Entry{
				&tux.Handle{c.name, makeSayHandler(c.say), nil},
			}, nil}

			muxTux := tux.ServeMux(group)

			mux := http.NewServeMux()

			for muxTux.Next() {
				pattern := muxTux.Pattern()
				handler := muxTux.Handler()

				fmt.Println(pattern)

				mux.Handle(pattern, handler)
			}
		})
	}
}

func makeSayHandler(say string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(say))
	})
}
