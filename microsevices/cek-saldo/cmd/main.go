package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"test_jajanomic/microservices/cek-saldo/models"
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
	r.HandleFunc("/api/saldo", HandleGetSaldo(db)).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("start servet at ", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func HandleGetSaldo(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req models.RekeningRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.RekeningResponse{
				IsError: true,
				Message: err.Error(),
			})
			return
		}

		rekening := models.Rekening{}
		if err := db.Model(rekening).Select("norek", "gold_balance").Where("norek = ?", req.Norek).First(&rekening).Error; err != nil {
			code := http.StatusInternalServerError
			if err != nil {
				code = http.StatusNotFound
			}
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(models.RekeningResponse{
				IsError: true,
				Message: err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(models.RekeningResponse{
			IsError: false,
			Data: struct {
				Norek       string  "json:\"norek,omitempty\""
				GoldBalance float32 "json:\"saldo,omitempty\""
			}{
				Norek:       rekening.Norek,
				GoldBalance: rekening.GoldBalance,
			},
		})
	}
}
