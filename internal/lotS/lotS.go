package lotS

import (
	"courseproject/internal/userS"
	"time"
)

type Lot struct {
	ID          int       `json:"id" gorm:"AUTO_INCREMENT; PRIMARY_KEY"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" `
	BuyPrice    float64   `json:"buy_price" `
	MinPrice    float64   `json:"min_price" gorm:"not null"`
	PriceStep   float64   `json:"price_step" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null"`
	EndAt       time.Time `json:"end_at"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatorID   int       `json:"creator" gorm:"not null"`
	BuyerID     int       `json:"buyer"`
}

type LotForJSON struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	BuyPrice    float64         `json:"buy_price"`
	MinPrice    float64         `json:"min_price"`
	PriceStep   float64         `json:"price_step"`
	Status      string          `json:"status"`
	EndAt       time.Time       `json:"end_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatorID   userS.ShortUser `json:"creator"`
	BuyerID     userS.ShortUser `json:"buyer"`
}

type LotTCU struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	MinPrice    float64   `json:"min_price"`
	PriceStep   float64   `json:"price_step"`
	EndAt       time.Time `json:"end_at"`
}

type Price struct {
	Price float64 `json:"price"`
}
