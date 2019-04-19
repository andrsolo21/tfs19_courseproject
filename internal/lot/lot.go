package lot

import (
	"courseproject/internal/lots"
	"courseproject/internal/storages"
	"errors"

	"time"
)

func Generate(l lots.LotTCU) lots.Lot {
	return lots.Lot{
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

func BuyLot(idUser int, idLot int, newPrice float64, db storages.INTT) (lots.Lot, error) {

	//проверка на ноль
	if newPrice <= 0 {
		return lots.Lot{}, errors.New("new price can't be less than 0")
	}

	/*us, err := db.Db().GetUserByID(idUser)
	if err != nil {
		return lotS.Lot{}, err
	}*/

	l, err := db.Db().GetLotByID(idLot)
	if err != nil {
		return lots.Lot{}, err
	}
	//проверка на минимльную цену
	if newPrice < l.MinPrice {
		return lots.Lot{}, errors.New("new price is less than the minimum")
	}

	//проверка на начальную цену
	//проверка на шаг
	if newPrice < l.BuyPrice+l.PriceStep && l.BuyPrice != 0 {
		return lots.Lot{}, errors.New("bad new price")
	}

	//проверка на тип лота
	switch l.Status {
	case "created":
		return lots.Lot{}, errors.New("lot exists, but is not traded")
	case "finished":
		return lots.Lot{}, errors.New("bidding on the lot is completed")
	}

	//проверка на время окончания действия лота
	if l.EndAt.Before(time.Now()) {
		return lots.Lot{}, errors.New("duration of the token has expired")
	}

	el, err := db.Db().SMBBuyIT(idUser, l, newPrice)

	return el, err
}
