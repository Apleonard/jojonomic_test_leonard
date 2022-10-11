package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"test_jajanomic/microservices/topup/models"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/teris-io/shortid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Print("Fail to load env")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	log.Println("connecting to db at ", dsn)

	conn, err := gorm.Open(postgres.Open(dsn), nil)
	if err != nil {
		log.Fatal("Error connect to db")
	}
	log.Println("connected to db at ", dsn)

	db = conn

	err = godotenv.Load()
	if err != nil {
		log.Print("Error load from file, read environemt from os environment")
	}

	servers := []string{os.Getenv("KAFKA_URL")}

	producer, err := sarama.NewSyncProducer(servers, nil)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/topup", HandleBuyback).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("start servet at ", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}

func HandleBuyback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req models.TopupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.TopupResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}
	rekening, err := getRekening(db, req.Norek)
	if err != nil {
		code := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound {
			code = http.StatusNotFound
		}

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(models.TopupResponse{
			IsError: true,
			Message: "rekening tidak ditemukan",
		})
		return
	}

	harga, err := getHarga(db, req.Amount)
	if err != nil {
		code := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound {
			code = http.StatusNotFound
		}

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(models.TopupResponse{
			IsError: true,
			Message: "harga tidak ditemukan",
		})
		return
	}

	reffID, err := shortid.Generate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.TopupResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	buybackParams := models.TopupParams{
		GoldWeight:         req.GoldWeight,
		Amount:             req.Amount,
		Norek:              req.Norek,
		ReffID:             reffID,
		HargaTopup:         harga.HargaTopup,
		HargaBuyback:       harga.HargaBuyback,
		CurrentGoldBalance: rekening.GoldBalance,
	}

	_, err = json.Marshal(&buybackParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.TopupResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	err = ProduceMessage(buybackParams)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(models.TopupResponse{
		IsError: false,
		ReffID:  reffID,
	})

}

func getRekening(db *gorm.DB, norek string) (*models.Rekening, error) {
	rekening := models.Rekening{}
	if err := db.Model(rekening).Where("norek = ?", norek).First(&rekening).Error; err != nil {
		return nil, err
	}

	return &rekening, nil
}

func getHarga(db *gorm.DB, buybakcAmount float64) (*models.Harga, error) {
	harga := models.Harga{}
	if err := db.Model(harga).Where("harga_topup = ?", buybakcAmount).Order("created_at desc").First(&harga).Error; err != nil {
		return nil, err
	}

	return &harga, nil
}

func ProduceMessage(params models.TopupParams) error {
	reqData, err := json.Marshal(&params)
	if err != nil {
		log.Fatal(err)
	}

	reqString := string(reqData)
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{os.Getenv("KAFKA_URL")}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	strTime := strconv.Itoa(int(time.Now().Unix()))

	msg := &sarama.ProducerMessage{
		Topic: os.Getenv("KAFKA_TOPIC"),
		Key:   sarama.StringEncoder(strTime),
		Value: sarama.StringEncoder(reqString),
	}

	producer.Input() <- msg

	return nil

}
