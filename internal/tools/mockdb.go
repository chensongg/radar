package tools

import "time"

type mockDB struct{}

var mockLoginDetails = map[string]LoginDetails{
	"alex": {
		AuthToken: "asdf",
		Username:  "alex",
	},
}

var mockCoinDetails = map[string]CoinDetails{
	"alex": {
		Coins:    142,
		Username: "alex",
	},
}

func (d *mockDB) GetUserLoginDetails(username string) *LoginDetails {
	time.Sleep(time.Millisecond * 50)
	var clientData, ok = mockLoginDetails[username]
	if !ok {
		return nil
	}
	return &clientData
}

func (d *mockDB) GetUserCoins(username string) *CoinDetails {
	time.Sleep(time.Millisecond * 50)
	var coinData, ok = mockCoinDetails[username]
	if !ok {
		return nil
	}
	return &coinData
}

func (d *mockDB) SetupDatabase() error {
	return nil
}
