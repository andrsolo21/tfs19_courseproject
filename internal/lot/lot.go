package lot

import (
	"courseproject/internal/lots"
	"courseproject/internal/storages"
	"errors"
	"time"
)

const (
	cr string = "created"
)

func Generate(l lots.LotTCU) (lots.Lot, error) {

	l, err := CheckLot(l)
	if err != nil {
		return lots.Lot{}, err
	}

	return lots.Lot{
		//ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		//BuyPrice:    l.BuyPrice,
		MinPrice:  l.MinPrice,
		PriceStep: l.PriceStep,
		Status:    l.Status,
		EndAt:     l.EndAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		//CreatorID:   user1.ToShort(),
		//BuyerID:     user2.ToShort(),
	}, nil
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
	if newPrice < l.BuyPrice+float64(l.PriceStep) && l.BuyPrice != 0 {
		return lots.Lot{}, errors.New("bad new price")
	}

	//проверка на тип лота
	switch l.Status {
	case cr:
		return lots.Lot{}, errors.New("lot exists, but is not traded")
	case "finished":
		return lots.Lot{}, errors.New("bidding on the lot is completed")
	}

	//проверка на время окончания действия лота
	if l.EndAt.Before(time.Now()) {
		return lots.Lot{}, errors.New("duration of the token has expired")
	}

	if l.CreatorID == idUser {
		return lots.Lot{}, errors.New("you cannot buy your lot")
	}

	el, err := db.SMBBuyIT(idUser, l, newPrice)

	return el, err
}

func CheckLot(l lots.LotTCU) (lots.LotTCU, error) {

	if l.Status == "" {
		l.Status = cr
	}

	if !(l.Status == cr || l.Status == "active" || l.Status == "finished") {
		return lots.LotTCU{}, errors.New("unexpected status")
	}

	if l.Title == "" {
		return lots.LotTCU{}, errors.New("title is empty")
	}

	if l.MinPrice < 1 {
		return lots.LotTCU{}, errors.New("bad min price")
	}

	if l.PriceStep < 1 {
		l.PriceStep = 1
	}

	return l, nil
}

func UpdateLot(userID int, l lots.LotTCU, id int, db storages.INTT) (lots.Lot, error) {

	lot, err := db.GetLotByID(id)
	if err != nil {
		return lots.Lot{}, errors.New("this lot doesn't exist")
	}

	switch l.Status {
	case "active":
		return lots.Lot{}, errors.New("lot is already trading")
	case "finished":
		return lots.Lot{}, errors.New("lot trading time is out")
	}

	if userID != lot.CreatorID {
		return lots.Lot{}, errors.New("trying to update someone else's lot")
	}

	l, err = CheckLot(l)
	if err != nil {
		return lots.Lot{}, err
	}

	lt := lots.Lot{
		Title:       l.Title,
		Description: l.Description,
		MinPrice:    l.MinPrice,
		PriceStep:   l.PriceStep,
		Status:      l.Status,
		EndAt:       l.EndAt,
		UpdatedAt:   time.Now(),
	}

	lt = db.UpdateLot(lt, id)

	return lt, nil
}

func DeleteLot(userID int, idLot int, db storages.INTT) error {

	lot, err := db.GetLotByID(idLot)
	if err != nil {
		return errors.New("this lot doesn't exist")
	}

	if userID != lot.CreatorID {
		return errors.New("trying to delete someone else's lot")
	}

	if lot.Status == cr {
		db.DeleteCrLor(idLot)
	} else {

		if !db.CheckDelLot(idLot) {
			return errors.New("lot already deleted")
		}
		lot.DeletedAt = time.Now()
		lot.UpdatedAt = time.Now()
		lot = db.UpdateLot(lot, idLot)
	}

	return nil
}