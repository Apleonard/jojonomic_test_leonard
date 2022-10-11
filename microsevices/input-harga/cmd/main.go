package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"test_jajanomic/microservices/input-harga/models"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/teris-io/shortid"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Fail to load env")
	}

	servers := []string{os.Getenv("KAFKA_URL")}

	producer, err := sarama.NewSyncProducer(servers, nil)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/input-harga", HandlerInputHarga).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("start servet at ", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func HandlerInputHarga(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var req models.InputHargaRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.InputHargaResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	reffId, err := shortid.Generate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.InputHargaResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}
	req.ReffId = reffId

	_, err = json.Marshal(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.InputHargaResponse{
			IsError: true,
			Message: err.Error(),
		})
		return
	}

	err = ProduceMessage(req)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.InputHargaResponse{
			IsError: true,
			ReffId:  reffId,
			Message: "Kafka not ready",
		})
		return

	}

	json.NewEncoder(w).Encode(models.InputHargaResponse{
		IsError: false,
		ReffId:  reffId,
	})

}

func ProduceMessage(params models.InputHargaRequest) error {
	reqData, err := json.Marshal(&params)
	if err != nil {
		log.Fatal(err)
	}

	reqString := string(reqData)
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{"KAFKA_URL"}
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
