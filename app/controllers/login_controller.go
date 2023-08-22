package controllers

import (
	"live_attendance/main/app/models"
	"net/http"

	"github.com/unrolled/render"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout: "layout",
	})

	_ = render.HTML(w, http.StatusOK, "login", map[string]interface{}{
		"title":   "Login",
		"success": GetFlash(w, r, "success"),
		"error":   GetFlash(w, r, "error"),
	})
}

func (server *Server) DoLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	employeeModel := models.Employee{}
	employee, err := employeeModel.FindByEmail(server.DB, email)

	if err != nil {
		SetFlash(w, r, "error", "Email or password invalid")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if employee.Password != password {
		SetFlash(w, r, "error", "Email or password invalid")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, sessionEmployee)
	session.Values["id"] = employee.ID
	session.Save(r, w)

	SetFlash(w, r, "success", "Login success")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionEmployee)

	session.Values["id"] = nil
	session.Save(r, w)

	SetFlash(w, r, "success", "Logout success")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
