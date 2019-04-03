package lot

import "time"

type Lot struct {
	ID          int       `json:"id"`
	Creator_id  int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"` //
	Min_price   float64   `json:"min_price"`
	Price_step  float64   `json:"price_step"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}
