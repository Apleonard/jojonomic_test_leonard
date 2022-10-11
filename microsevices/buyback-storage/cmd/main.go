package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"test_jajanomic/microservices/buyback-storage/models"
	"time"

	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Print("Error load from file, read environemt from os environment")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	conn, err := gorm.Open(postgres.Open(dsn), nil)
	if err != nil {
		log.Fatal("Error connect to db")
	}

	db = conn

	ConsumeMessage()

}

func ConsumeMessage() {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Consumer.Return.Errors = true
	servers := []string{os.Getenv("KAFKA_URL")}

	consumer, err := sarama.NewConsumer(servers, nil)
	if err != nil {
		log.Println(err)
	}

	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(os.Getenv("KAFKA_TOPIC"), int32(0), sarama.OffsetNewest)
	if err != nil {
		log.Println(err)
	}

	defer partitionConsumer.Close()

	var harga = models.Transaction{}
	for {
		select {

		case err := <-partitionConsumer.Errors():
			log.Println(err)
			break

		case msg := <-partitionConsumer.Messages():

			err := json.Unmarshal(msg.Value, &harga)
			if err != nil {
				log.Fatal(err)
			}

			var buybackParams models.BuybackParams
			err = json.Unmarshal(msg.Value, &buybackParams)
			if err != nil {
				return
			}

			goldBalance := buybackParams.CurrentGoldBalance - buybackParams.GoldWeight
			transaction := models.Transaction{
				ReffID:       buybackParams.ReffID,
				Norek:        buybackParams.Norek,
				Type:         "buyback",
				GoldWeight:   buybackParams.GoldWeight,
				GoldBalance:  goldBalance,
				HargaTopup:   buybackParams.HargaTopup,
				HargaBuyback: buybackParams.HargaBuyback,
				CreatedAt:    int(time.Now().Unix()),
			}

			if err := db.Create(&transaction).Error; err != nil {
				db.Rollback()
			}

			if err := db.Model(models.Rekening{}).Where("norek = ?", buybackParams.Norek).Update("gold_balance", goldBalance).Error; err != nil {
				db.Rollback()
			}

			log.Printf("success save data : %s", string(msg.Value))
			continue
		}
	}
}
