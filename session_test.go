package scs

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/porcupineyhairs/scs/rqlitestore"
	"github.com/rqlite/gorqlite"
)

func TestPoc(t *testing.T) {
	// Establish connection to rqlite.
	conn, err := gorqlite.Open("http://localhost:4001/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Initialize a new session manager and configure it to use rqlitestore as the session store.
	sessionManager := New()
	sessionManager.Store = rqlitestore.New(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		sessionManager.Put(r.Context(), "message", "Hello from a session!")
	})
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		msg := sessionManager.GetString(r.Context(), "message")
		io.WriteString(w, msg)
	})

	http.ListenAndServe(":4000", sessionManager.LoadAndSave(mux))
}
