package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello world")
	d, err := io.ReadAll(r.Body)

	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte("oh Snap.."))

		http.Error(rw, "oh snap..", http.StatusBadRequest)
		return
	}
	h.l.Printf("data => %s", d)
	fmt.Fprintf(rw, "Hello %s", d)
}
