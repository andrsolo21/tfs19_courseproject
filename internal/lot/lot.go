package lot

import (
	"errors"
	"time"

	"gitlab.com/andrsolo21/courseproject/internal/lots"
	"gitlab.com/andrsolo21/courseproject/internal/storages"
)

const (
	cr string = "created"
	fn string = "finished"
	ac string = "active"
)

func Generate(l lots.LotTCU) (lots.Lot, error) {

	l, err := CheckLot(l)
	if err != nil {
		return lots.Lot{}, err
	}

	return lots.Lot{
		//ID:          l.ID,
		Title:       lots.B(l.Title),
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

func CheckLot(l lots.LotTCU) (lots.LotTCU, error) {

	if l.Status == "" {
		l.Status = cr
	}

	if !(l.Status == cr || l.Status == ac || l.Status == fn) {
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

	err = CheckLotForTraiding(newPrice, l, idUser)
	if err != nil {
		return l, err
	}

	el, err := db.SMBBuyIT(idUser, l, newPrice)

	return el, err
}

func CheckLotForTraiding(newPrice float64, l lots.Lot, idUser int) error {

	//проверка на минимльную цену
	if newPrice < l.MinPrice {
		return errors.New("new price is less than the minimum")
	}

	//проверка на начальную цену
	//проверка на шаг
	if newPrice < l.BuyPrice+float64(l.PriceStep) && l.BuyPrice != 0 {
		return errors.New("bad new price")
	}

	//проверка на тип лота
	switch l.Status {
	case cr:
		return errors.New("lot exists, but is not traded")
	case fn:
		return errors.New("bidding on the lot is completed")
	}

	//проверка на время окончания действия лота
	if l.EndAt.Before(time.Now()) {
		return errors.New("duration of the token has expired")
	}

	if l.CreatorID == idUser {
		return errors.New("you cannot buy your lot")
	}

	if l.BuyerID == idUser {
		return errors.New("you already buy this lot")
	}
	return nil
}

func UpdateLot(userID int, l lots.LotTCU, id int, db storages.INTT) (lots.Lot, error) {

	lt, err := db.GetLotByID(id)
	if err != nil {
		return lots.Lot{}, errors.New("this lot doesn't exist")
	}

	l, err = CheckLotForUpdate(userID, lt, l)
	if err != nil {
		return lots.Lot{}, err
	}

	lt2 := lots.Lot{
		Title:       l.Title,
		Description: l.Description,
		MinPrice:    l.MinPrice,
		PriceStep:   l.PriceStep,
		Status:      l.Status,
		EndAt:       l.EndAt,
		UpdatedAt:   time.Now(),
	}

	lt2 = db.UpdateLot(lt2, id)

	return lt2, nil
}

func CheckLotForUpdate(userID int, lt lots.Lot, l lots.LotTCU) (lots.LotTCU, error) {
	if userID != lt.CreatorID {
		return lots.LotTCU{}, errors.New("trying to update someone else's lot")
	}

	switch lt.Status {
	case ac:
		return lots.LotTCU{}, errors.New("lot is already trading")
	case fn:
		return lots.LotTCU{}, errors.New("lot trading time is out")
	}

	l, err := CheckLot(l)

	return l, err
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
		_, el := db.CheckDelLot(idLot)
		if el.DeletedAt.After(time.Now().Add(time.Hour * 9000 * -1)) {
			return errors.New("lot already deleted")
		}

		lot.DeletedAt = time.Now()
		lot.UpdatedAt = time.Now()
		lot.Status = "finished"
		lot = db.UpdateLot(lot, idLot)
	}

	return nil
}

func Separate(lts []lots.Lot) []lots.Lot {
	var lts2 []lots.Lot
	for _, el := range lts {
		if el.DeletedAt.Before(time.Now().Add(time.Hour * 9000 * -1)) {
			lts2 = append(lts2, el)
		}
	}
	return lts2
}
