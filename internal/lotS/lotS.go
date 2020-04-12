package lots

import (
	"time"

	"gitlab.com/andrsolo21/courseproject/internal/users"
)

type Lot struct {
	ID          int       `json:"id" gorm:"AUTO_INCREMENT; PRIMARY_KEY"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" `
	BuyPrice    float64   `json:"buy_price" `
	MinPrice    float64   `json:"min_price" gorm:"not null"`
	PriceStep   int       `json:"price_step" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null"`
	EndAt       time.Time `json:"end_at"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatorID   int       `json:"creator" gorm:"not null"`
	BuyerID     int       `json:"buyer,omitempty"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
}

type LotForJSON struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	BuyPrice    float64         `json:"buy_price,omitempty"`
	MinPrice    float64         `json:"min_price"`
	PriceStep   int             `json:"price_step"`
	Status      string          `json:"status"`
	EndAt       time.Time       `json:"end_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatorID   users.ShortUser `json:"creator"`
	BuyerID     users.ShortUser `json:"buyer,omitempty"`
	DeletedAt   time.Time       `json:"deleted_at,omitempty"`
}

type LotTCU struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	MinPrice    float64   `json:"min_price"`
	PriceStep   int       `json:"price_step"`
	EndAt       time.Time `json:"end_at"`
	Status      string    `json:"status"`
}

type Price struct {
	Price float64 `json:"price"`
}

func B(a string) string {
	return a
}
