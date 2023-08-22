package controllers

import (
	"net/http"

	"github.com/unrolled/render"
)

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout: "layout",
	})

	_ = render.HTML(w, http.StatusOK, "register", map[string]interface{}{
		"title":   "Register",
		"success": GetFlash(w, r, "success"),
		"error":   GetFlash(w, r, "error"),
	})
}
