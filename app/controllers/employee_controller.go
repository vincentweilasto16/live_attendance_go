package controllers

import (
	"live_attendance/main/app/models"
	"live_attendance/main/app/utils"
	"net/http"

	"github.com/google/uuid"

	"io"
	"os"
	"path/filepath"
)

func (server *Server) AddEmployee(w http.ResponseWriter, r *http.Request) {
	var employee models.Employee

	if err := r.ParseMultipartForm(10 << 20); // 10 MB limit for image size
	err != nil {
		SetFlash(w, r, "error", "Invalid data")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}

	employee.ID = uuid.New().String()
	employee.Name = r.FormValue("name")
	employee.Email = r.FormValue("email")
	employee.Password = r.FormValue("password")

	// Get the file from the form data
	imageFile, handler, err := r.FormFile("image")
	if err != nil {

		SetFlash(w, r, "error", "Failed to upload image")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	defer imageFile.Close()

	newFileName := employee.ID + "_" + handler.Filename

	// Determine the file path to save the image in the 'assets' folder
	savePath := filepath.Join("assets/user_image", newFileName)

	// Create a new file on the server
	newFile, err := os.Create(savePath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to create file on server")
		return
	}
	defer newFile.Close()

	// Copy the uploaded file to the new file on the server
	_, err = io.Copy(newFile, imageFile)
	if err != nil {
		http.Error(w, "Unable to copy file", http.StatusInternalServerError)
		return
	}

	employee.Image = newFileName

	result := server.DB.Create(&employee)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add employee")
		return
	}

	SetFlash(w, r, "success", "Register success, please login with your account")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
