package lot

import (
	"courseproject/internal/lots"
	"testing"
)

type testpair struct {
	value lots.LotTCU
	out   bool
}

func TestCheckLot(t *testing.T) {

	var tests = []testpair{
		{lots.LotTCU{
			Status: "created",
			MinPrice: 1,
			Title: "smth",
		}, true},
		{lots.LotTCU{
			Status: "",
			MinPrice: 2,
			Title: "smth",
		}, true},
		{lots.LotTCU{
			Status: "finished",
			MinPrice: 0,
			Title: "smth",
		}, false},
		{lots.LotTCU{
			Status: "finished",
			MinPrice: 5,
			Title: "",
		}, false},
		{lots.LotTCU{
			Status: "badStatus",
			MinPrice: 7,
			Title: "smth",
		}, false},
	}

	for _, pair := range tests {
		_, err := CheckLot(pair.value)
		if (err == nil) != pair.out {
			t.Error(
				"For", pair.value,
				"expected", pair.out,
				"got", pair.out,
			)
		}
	}
}
