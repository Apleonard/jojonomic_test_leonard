package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"test_jajanomic/microservices/buyback-storage/models"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
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

	r := getKafkaReader(os.Getenv("KAFKA_URL"), os.Getenv("KAFKA_TOPIC"), os.Getenv("KAFKA_GROUP_ID"))
	ctx := context.Background()
	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}
		log.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Key))

		if err := SaveBuyback(db, m.Value); err != nil {
			log.Println(err.Error())
			continue
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}
}

func SaveBuyback(db *gorm.DB, data []byte) error {
	var buybackParams models.BuybackParams
	if err := json.Unmarshal(data, &buybackParams); err != nil {
		return fmt.Errorf("unmarshall data error : %s", err.Error())
	}

	conn := db.Begin()

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

	if err := conn.Create(&transaction).Error; err != nil {
		conn.Rollback()
		return err
	}

	if err := conn.Model(models.Rekening{}).Where("norek = ?", buybackParams.Norek).Update("gold_balance", goldBalance).Error; err != nil {
		conn.Rollback()
		return err
	}

	return conn.Commit().Error
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})
}
