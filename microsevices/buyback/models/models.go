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

type Harga struct {
	AdminId      string    `json:"admin_id"`
	ReffId       string    `json:"reff_id"`
	HargaTopup   float64   `json:"harga_topup"`
	HargaBuyback float64   `json:"harga_buyback"`
	CreatedAt    time.Time `json:"created_at"`
}

func (Harga) TableName() string {
	return "tbl_harga"
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

type BuybackRequest struct {
	GoldWeight float32 `json:"gram"`
	Amount     float64 `json:"harga"`
	Norek      string  `json:"norek"`
}

type BuybackResponse struct {
	IsError bool   `json:"error"`
	ReffID  string `json:"reff_id,omitempty"`
	Message string `json:"message,omitempty"`
}
