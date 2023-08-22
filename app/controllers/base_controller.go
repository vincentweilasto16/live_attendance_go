package controllers

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"

	"live_attendance/main/app/models"

	"github.com/corona10/goimagehash"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nfnt/resize"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBDriver   string
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var sessionFlash = "flash-session"
var sessionEmployee = "employee-session"

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcome to " + appConfig.AppName)

	server.initializeDB(dbConfig)
	server.initializeRoutes()
	// seeders.DBSeed(server.DB)
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *Server) initializeDB(dbConfig DBConfig) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
	server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed on connecting to the database server")
	}

	server.dbMigrate()
}

func (server *Server) dbMigrate() {
	for _, model := range models.RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal()
		}
	}

	fmt.Println("Database migrated successfully")
}

func SetFlash(w http.ResponseWriter, r *http.Request, name string, value string) {
	session, err := store.Get(r, sessionFlash)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.AddFlash(value, name)
	session.Save(r, w)
}

func GetFlash(w http.ResponseWriter, r *http.Request, name string) []string {
	session, err := store.Get(r, sessionFlash)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	fm := session.Flashes(name)
	if len(fm) < 0 {
		return nil
	}

	session.Save(r, w)
	var flashes []string
	for _, fl := range fm {
		flashes = append(flashes, fl.(string))
	}

	return flashes
}

func IsLoggedIn(r *http.Request) bool {
	session, _ := store.Get(r, sessionEmployee)
	if session.Values["id"] == nil {
		return false
	}

	return true
}

func compareImages(imagePath1, imagePath2 string) float64 {
	image1, err := loadImage(imagePath1)
	if err != nil {
		log.Fatal(err)
	}

	image2, err := loadImage(imagePath2)
	if err != nil {
		log.Fatal(err)
	}

	// Resize images to the same dimensions
	width := 8
	height := 8
	resizedImg1 := resize.Resize(uint(width), uint(height), image1, resize.Lanczos3)
	resizedImg2 := resize.Resize(uint(width), uint(height), image2, resize.Lanczos3)

	// Calculate image hashes
	hash1 := calculateImageHash(resizedImg1)
	hash2 := calculateImageHash(resizedImg2)

	return calculateImageSimilarity(hash1, hash2)
}

func loadImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func calculateImageHash(img image.Image) *goimagehash.ImageHash {
	hash, _ := goimagehash.AverageHash(img)
	return hash
}

func calculateImageSimilarity(hash1, hash2 *goimagehash.ImageHash) float64 {
	hashBits := hash1.Bits()
	maxDistance := float64(hashBits)

	distance, _ := hash1.Distance(hash2)
	normalizedSimilarity := 1.0 - (float64(distance) / maxDistance)
	return normalizedSimilarity
}
