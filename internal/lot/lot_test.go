package lot

import (
	"testing"
	"time"

	"gitlab.com/andrsolo21/courseproject/internal/lots"
)

type testpair struct {
	value lots.LotTCU
	out   bool
}

type testpair2 struct {
	value    lots.Lot
	out      bool
	buyPrice float64
	buyerID  int
}

type testpair3 struct {
	lt     lots.Lot
	ltTCU  lots.LotTCU
	userID int
	out    bool
}

func TestCheckLot(t *testing.T) {

	var tests = []testpair{
		{lots.LotTCU{
			Status:   "created",
			MinPrice: 1,
			Title:    "smth",
		}, true},
		{lots.LotTCU{
			Status:    "",
			MinPrice:  2,
			Title:     "smth",
			PriceStep: 0,
		}, true},
		{lots.LotTCU{
			Status:   "finished",
			MinPrice: 0,
			Title:    "smth",
		}, false},
		{lots.LotTCU{
			Status:   "finished",
			MinPrice: 5,
			Title:    "",
		}, false},
		{lots.LotTCU{
			Status:    "badStatus",
			MinPrice:  7,
			Title:     "smth",
			PriceStep: 0,
		}, false},
	}

	for i, pair := range tests {
		lr, err := CheckLot(pair.value)
		if (err == nil) != pair.out && lr.PriceStep < 1 {
			t.Error(
				//"For", pair.value,
				"expected", pair.out,
				"got", pair.out,
				"test ID: ", i,
			)
		}
	}
}

func TestCheckLotForTraiding(t *testing.T) {

	var tests = []testpair2{
		{lots.Lot{Status: "active",
			MinPrice:  1,
			PriceStep: 1,
			Title:     "smth",
			BuyPrice:  5,
			EndAt:     time.Now().Add(time.Hour),
			CreatorID: 5,
			BuyerID:   9}, true, 10, 7},
		{lots.Lot{Status: "active",
			MinPrice:  200,
			PriceStep: 1,
			Title:     "smth",
			BuyPrice:  5,
			EndAt:     time.Now().Add(time.Hour),
			CreatorID: 5,
			BuyerID:   9}, false, 10, 7},
		{lots.Lot{Status: "active",
			MinPrice:  1,
			PriceStep: 200,
			Title:     "smth",
			BuyPrice:  5,
			EndAt:     time.Now().Add(time.Hour),
			CreatorID: 5,
			BuyerID:   9}, false, 20, 7},
		{lots.Lot{Status: "created",
			MinPrice:  1,
			PriceStep: 1,
			Title:     "smth",
			BuyPrice:  5,
			EndAt:     time.Now().Add(time.Hour),
			CreatorID: 5,
			BuyerID:   9}, false, 20, 7},
		{lots.Lot{Status: "finished",
			MinPrice:  1,
			PriceStep: 1,
			Title:     "smth",
			BuyPrice:  5,
			EndAt:     time.Now().Add(time.Hour),
			CreatorID: 5,
			BuyerID:   9}, false, 20, 7},
		{lots.Lot{Status: "active",
			MinPrice:  1,
			PriceStep: 1,
			Title:     "smth",
			BuyPrice:  5,
			EndAt:     time.Now().Add(time.Hour),
			CreatorID: 5,
			BuyerID:   9}, false, 20, 5},
		{lots.Lot{Status: "active",
			MinPrice:  1,
			PriceStep: 1,
			Title:     "smth",
			BuyPrice:  5,
			EndAt:     time.Now().Add(time.Hour),
			CreatorID: 5,
			BuyerID:   9}, false, 20, 9},
		{lots.Lot{Status: "active",
			MinPrice:  1,
			PriceStep: 1,
			Title:     "smth",
			BuyPrice:  5,
			EndAt:     time.Now().Add(-1 * time.Hour),
			CreatorID: 5,
			BuyerID:   9}, false, 20, 7},
	}

	for i, pair := range tests {
		err := CheckLotForTraiding(pair.buyPrice, pair.value, pair.buyerID)
		if (err == nil) != pair.out {
			t.Error(
				//"For", pair.value,
				"expected", pair.out,
				"got", pair.out,
				"test ID: ", i,
			)
		}
	}
}

func TestCheckLotForUpdate(t *testing.T) {

	var tests = []testpair3{
		{lots.Lot{Status: "created",
			CreatorID: 5,
		},
			lots.LotTCU{
				Status:   "created",
				MinPrice: 1,
				Title:    "smth",
			}, 5, true},
		{lots.Lot{Status: "active",
			CreatorID: 5,
		},
			lots.LotTCU{
				Status:   "created",
				MinPrice: 1,
				Title:    "smth",
			}, 5, false},
		{lots.Lot{Status: "finished",
			CreatorID: 5,
		},
			lots.LotTCU{
				Status:   "created",
				MinPrice: 1,
				Title:    "smth",
			}, 5, false},
		{lots.Lot{Status: "created",
			CreatorID: 5,
		},
			lots.LotTCU{
				Status:   "active",
				MinPrice: 1,
				Title:    "smth",
			}, 5, true},
		{lots.Lot{Status: "created",
			CreatorID: 5,
		},
			lots.LotTCU{
				Status:   "active",
				MinPrice: 1,
				Title:    "smth",
			}, 7, false},
		{lots.Lot{Status: "created",
			CreatorID: 5,
		},
			lots.LotTCU{
				Status:   "active",
				MinPrice: 1,
				Title:    "",
			}, 5, false},
		{lots.Lot{Status: "created",
			CreatorID: 5,
		},
			lots.LotTCU{
				Status:   "finished",
				MinPrice: 1,
				Title:    "",
			}, 5, true},
	}

	for i, pair := range tests {
		l, err := CheckLotForUpdate(pair.userID, pair.lt, pair.ltTCU)
		if (err == nil) != pair.out && (l.Status == pair.ltTCU.Status || pair.ltTCU.Status == "") {
			t.Error(
				//"For", pair.value,
				"expected", pair.out,
				"got", pair.out,
				"test ID: ", i,
			)
		}
	}
}
