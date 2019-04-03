package lot

import (
	"courseproject/internal/user"
	"time"
)

type Lot struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	BuyPrice    float64   `json:"buy_price"`
	MinPrice    float64   `json:"min_price"`
	PriceStep   float64   `json:"price_step"`
	Status      string    `json:"status"`
	EndAt       time.Time `json:"end_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatorID   int       `json:"creator"`
	BuyerID     int       `json:"buyer"`
}

type LotForJSON struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	BuyPrice    float64        `json:"buy_price"`
	MinPrice    float64        `json:"min_price"`
	PriceStep   float64        `json:"price_step"`
	Status      string         `json:"status"`
	EndAt       time.Time      `json:"end_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatorID   user.ShortUser `json:"creator"`
	BuyerID     user.ShortUser `json:"buyer"`
}

type LotTCU struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	MinPrice    float64   `json:"min_price"`
	PriceStep   float64   `json:"price_step"`
	EndAt       time.Time `json:"end_at"`
}

func (l LotTCU) Generate()Lot{
	return Lot{
		//ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		//BuyPrice:    l.BuyPrice,
		MinPrice:    l.MinPrice,
		PriceStep:   l.PriceStep,
		Status:      "created",
		EndAt:       l.EndAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		//CreatorID:   user1.ToShort(),
		//BuyerID:     user2.ToShort(),
	}
}
