package lot

import (
	"courseproject/internal/lotS"

	"time"
)

func Generate(l lotS.LotTCU) lotS.Lot {
	return lotS.Lot{
		//ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		//BuyPrice:    l.BuyPrice,
		MinPrice:  l.MinPrice,
		PriceStep: l.PriceStep,
		Status:    "created",
		EndAt:     l.EndAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		//CreatorID:   user1.ToShort(),
		//BuyerID:     user2.ToShort(),
	}
}
