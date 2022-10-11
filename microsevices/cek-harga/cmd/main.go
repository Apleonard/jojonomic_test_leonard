package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"test_jajanomic/microservices/cek-harga/models"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error load from file, read environemt from os environment")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), nil)
	if err != nil {
		log.Fatal("Error connect to db")
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/check-harga", HandleGetHarga(db)).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("start serve at ", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func HandleGetHarga(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var harga models.Harga
		if err := db.Model(harga).Order("created_at desc").First(&harga).Error; err != nil {
			code := http.StatusInternalServerError
			if err == gorm.ErrRecordNotFound {
				code = http.StatusNotFound
			}

			w.WriteHeader(code)
			json.NewEncoder(w).Encode(models.Response{
				IsError: true,
				Message: err.Error(),
			})
		}

		log.Printf("found harga : %+v", harga)
		json.NewEncoder(w).Encode(models.Response{
			IsError: false,
			Data: struct {
				HargaTopup   float64 "json:\"harga_topup\""
				HargaBuyback float64 "json:\"harga_buyback\""
			}{
				HargaTopup:   harga.HargaTopup,
				HargaBuyback: harga.HargaBuyback,
			},
		})
	}
}
