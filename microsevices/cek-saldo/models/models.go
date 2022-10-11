package models

import "time"

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

type RekeningRequest struct {
	Norek string `json:"norek"`
}

type RekeningResponse struct {
	IsError bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Norek       string  `json:"norek,omitempty"`
		GoldBalance float32 `json:"saldo,omitempty"`
	} `json:"data,omitempty"`
}
