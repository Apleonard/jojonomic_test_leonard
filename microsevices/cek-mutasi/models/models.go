package models

type TransactionRequest struct {
	Norek     string `json:"norek"`
	StartDate int32  `json:"start_date"`
	EndDate   int32  `json:"end_date"`
}

type TransactionResponseItem struct {
	CreatedAt    int32   `json:"date"`
	Type         string  `json:"type"`
	GoldWeight   float32 `json:"gram"`
	HargaTopup   float64 `json:"harga_topup"`
	HargaBuyback float64 `json:"harga_buyback"`
	GoldBalance  float32 `json:"saldo"`
}

type TransactionResponse struct {
	IsError bool                       `json:"error"`
	Data    []*TransactionResponseItem `json:"data,omitempty"`
	Message string                     `json:"message,omitempty"`
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

type ListTransaction []*Transaction

func (lt ListTransaction) ToResponseItems() []*TransactionResponseItem {
	list := make([]*TransactionResponseItem, len(lt))
	for i, v := range lt {
		list[i] = &TransactionResponseItem{
			CreatedAt:    int32(v.CreatedAt),
			Type:         v.Type,
			GoldWeight:   v.GoldWeight,
			GoldBalance:  v.GoldBalance,
			HargaTopup:   v.HargaTopup,
			HargaBuyback: v.HargaBuyback,
		}
	}

	return list
}
