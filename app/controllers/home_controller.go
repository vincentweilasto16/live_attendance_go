package controllers

import (
	"fmt"
	"io"
	"live_attendance/main/app/models"
	"live_attendance/main/app/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/unrolled/render"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {

	render := render.New(render.Options{
		Layout: "layout",
	})

	if !IsLoggedIn(r) {
		SetFlash(w, r, "error", "You need to login first!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	session, _ := store.Get(r, "employee-session")

	employeeModel := models.Employee{}
	employee, _ := employeeModel.FindByID(server.DB, session.Values["id"].(string))

	_ = render.HTML(w, http.StatusOK, "home", map[string]interface{}{
		"title":    "Home",
		"success":  GetFlash(w, r, "success"),
		"error":    GetFlash(w, r, "error"),
		"employee": employee,
	})
}

func (server *Server) ClockIn(w http.ResponseWriter, r *http.Request) {
	var attendance models.Attendance

	if err := r.ParseMultipartForm(10 << 20); // 10 MB limit for image size
	err != nil {
		SetFlash(w, r, "error", "Invalid data")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	session, _ := store.Get(r, "employee-session")
	employeeModel := models.Employee{}
	employee, _ := employeeModel.FindByID(server.DB, session.Values["id"].(string))

	_, err := attendance.FindByIDAndClockIn(server.DB, employee.ID)
	if err == nil {
		SetFlash(w, r, "error", "You already clock-in today!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get the file from the form data
	imageFile, _, err := r.FormFile("clockin-image")
	if err != nil {
		SetFlash(w, r, "error", "Failed to upload image")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	defer imageFile.Close()

	attendance.ID = uuid.New().String()
	newFileName := "clock-in-image" + "_" + attendance.ID + ".jpg"

	// Determine the file path to save the image in the 'assets' folder
	savePath := filepath.Join("assets/attendance_image", newFileName)

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

	storedImagePath := filepath.Join("assets/user_image/", employee.Image)

	similarity := compareImages(storedImagePath, savePath)
	fmt.Println(similarity)
	if similarity >= 0.8 {
		attendance.EmployeeID = employee.ID
		attendance.LastClockIn = time.Now()

		attendance.ClockInImage = newFileName

		result := server.DB.Create(&attendance)
		if result.Error != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add attendance record")
			return
		}

		SetFlash(w, r, "success", "Success Clock In")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		removePath := filepath.Join("assets/attendance_image", newFileName)
		newFile.Close()

		err := os.Remove(removePath)
		if err != nil {
			fmt.Printf("Error removing file: %v\n", err)
			return
		}

		SetFlash(w, r, "error", "Your picture are not match, please take your picture again!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func (server *Server) ClockOut(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(10 << 20); // 10 MB limit for image size
	err != nil {
		SetFlash(w, r, "error", "Invalid data")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	session, _ := store.Get(r, "employee-session")
	employeeModel := models.Employee{}
	employee, _ := employeeModel.FindByID(server.DB, session.Values["id"].(string))

	attendanceModel := models.Attendance{}
	_, err := attendanceModel.FindByIDAndClockOut(server.DB, employee.ID)
	if err == nil {
		SetFlash(w, r, "error", "You already clock-out today!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	attendance, err := attendanceModel.FindByIDAndClockIn(server.DB, employee.ID)

	// Get the file from the form data
	imageFile, _, err := r.FormFile("clockout-image")
	if err != nil {
		SetFlash(w, r, "error", "Failed to upload image")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	defer imageFile.Close()

	newFileName := "clock-out-image" + "_" + attendance.ID + ".jpg"

	// Determine the file path to save the image in the 'assets' folder
	savePath := filepath.Join("assets/attendance_image", newFileName)

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

	storedImagePath := filepath.Join("assets/user_image/", employee.Image)

	similarity := compareImages(storedImagePath, savePath)
	fmt.Println(similarity)
	if similarity >= 0.8 {
		result := server.DB.Model(&models.Attendance{}).Where("id = ?", attendance.ID).Updates(map[string]interface{}{
			"clock_out_image": newFileName,
			"last_clock_out":  time.Now(),
		})
		if result.Error != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update attendance record")
			return
		}

		SetFlash(w, r, "success", "Success Clock Out")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		removePath := filepath.Join("assets/attendance_image", newFileName)
		newFile.Close()

		err := os.Remove(removePath)
		if err != nil {
			fmt.Printf("Error removing file: %v\n", err)
			return
		}

		SetFlash(w, r, "error", "Your picture are not match, please take your picture again!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
