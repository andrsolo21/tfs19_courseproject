package lot

import (
	storages "courseproject/internal/database"
	"courseproject/internal/lotS"
	"errors"

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

func BuyLot(idUser int, idLot int, newPrice float64, db storages.INTT) (lotS.Lot, error) {

	//проверка на ноль
	if newPrice <= 0 {
		return lotS.Lot{}, errors.New("new price can't be less than 0")
	}

	/*us, err := db.Db().GetUserByID(idUser)
	if err != nil {
		return lotS.Lot{}, err
	}*/

	l, err := db.Db().GetLotByID(idLot)
	if err != nil {
		return lotS.Lot{}, err
	}
	//проверка на минимльную цену
	if newPrice < l.MinPrice {
		return lotS.Lot{}, errors.New("new price is less than the minimum")
	}

	//проверка на начальную цену
	//проверка на шаг
	if newPrice < l.BuyPrice+l.PriceStep && l.BuyPrice != 0 {
		return lotS.Lot{}, errors.New("bad new price")
	}

	//проверка на тип лота
	switch l.Status {
	case "created":
		return lotS.Lot{}, errors.New("lot exists, but is not traded")
	case "finished":
		return lotS.Lot{}, errors.New("bidding on the lot is completed")
	}

	//проверка на время окончания действия лота
	if l.EndAt.Before(time.Now()) {
		return lotS.Lot{}, errors.New("duration of the token has expired")
	}

	el, err := db.Db().SMBBuyIT(idUser, l, newPrice)

	return el, nil
}
