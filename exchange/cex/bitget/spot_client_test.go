package bitget

import "testing"

func TestClient_CoinsInfo(t *testing.T) {
	client := SetUp()
	coins, err := client.CoinsInfo("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(coins)
}

func TestClient_SpotSymbols(t *testing.T) {
	client := SetUp()
	coins, err := client.SpotSymbols("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(coins)
}
func TestClient_SpotTicker(t *testing.T) {
	client := SetUp()
	rsp, err := client.SpotTickers("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rsp)

}
