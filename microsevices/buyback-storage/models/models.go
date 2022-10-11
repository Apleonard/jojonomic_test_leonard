package models

import "time"

type BuybackParams struct {
	GoldWeight         float32 `json:"gram"`
	Amount             float64 `json:"harga"`
	Norek              string  `json:"norek"`
	ReffID             string  `json:"reff_id,omitempty"`
	HargaTopup         float64 `json:"harga_topup"`
	HargaBuyback       float64 `json:"harga_buyback"`
	CurrentGoldBalance float32 `json:"current_gold_balance"`
}

type Rekening struct {
	ReffID       string    `json:"reff_id"`
	Norek        string    `json:"norek"`
	CustomerName string    `json:"customer_name"`
	GoldBalance  float32   `json:"gold_balance"`
	CreateAt     time.Time `json:"created_at"`
}

func (Rekening) TableName() string {
	return "tbl_rekening"
}

type Transaction struct {
	ReffID       string  `json:"reff_id"`
	Norek        string  `json:"norek"`
	Type         string  `json:"type"`
	GoldWeight   float32 `json:"gold_weight"`
	HargaTopup   float64 `json:"harga_topup"`
	HargaBuyback float64 `json:"harga_buyback"`
	GoldBalance  float32 `json:"gold_balance"`
	CreatedAt    int     `json:"created_at"`
}

func (Transaction) TableName() string {
	return "tbl_transaksi"
}
