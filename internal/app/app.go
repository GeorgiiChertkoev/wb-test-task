package app

import (
	"net/http"
)

type App struct {
	a int
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Home!"))
}
