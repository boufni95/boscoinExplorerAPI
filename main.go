package main

import (
	"fmt"
	"math/big"
)

func main() {
	RetriveAndCalc()
	//GetAccount("GA4GYZGF4UPSDVY2FDFL7WA42IGEYXIPXLDLZ4QAFNPALRTTXAFULMIJ")
}
func RetriveAndCalc() {
	data := map[string]interface{}{
		"total":  0,
		"txAddr": "",
		"txNum":  0,
	}
	tmp := big.NewInt(0)
	var f map[string]interface{}
	var tot *big.Int
	if data["txAddr"] == "" {
		tot = big.NewInt(0)
		f = GetFrozenInter("/api/v1/frozen-accounts")
	} else {
		n := big.NewInt(data["total"].(int64))
		ten6 := big.NewInt(1000000)
		n.Mul(n, ten6)
		tot = n
		f = GetFrozenInter(data["txAddr"].(string))
	}
	/*var out bytes.Buffer
	json.Indent(&out, b, " ", "    ")
	SaveToFile("frozen.json", out.Bytes())*/
	var numberOf int
	var hrefNext string
	//for tmp.Sign() >= 0 {
	for i := 0; i < 5; i++ {
		hrefNext = GetNext(f)
		//hrefPrev := GetPrev(f)
		f = GetFrozenInter(hrefNext)
		tmp, numberOf = CountFreeze(f)

		tot.Add(tot, tmp)
	}
	fmt.Println("tot:", tot)
	fmt.Println("num:", numberOf)
	fmt.Println("ref:", hrefNext)
}
