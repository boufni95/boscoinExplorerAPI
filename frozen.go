package boscoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

type FrozenState string

const (
	frozen   FrozenState = "frozen"
	melting  FrozenState = "melting"
	returned FrozenState = "returned"
)

func GetFrozenBytes(moreP string) ([]byte, error) {

	var path strings.Builder
	path.WriteString("https://mainnet.blockchainos.org")
	path.WriteString(moreP)

	res, err := http.Get(path.String())
	if err != nil {
		log.Fatal(err)
	}
	received, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return received, nil
}
func GetFrozenInter(moreP string) (map[string]interface{}, error) {

	received, err := GetFrozenBytes(moreP)
	if err != nil {
		return nil, err
	}
	i := make(map[string]interface{})
	err = json.Unmarshal(received, &i)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func GetNext(i map[string]interface{}) (string, error) {
	//i := make(map[string]interface{})
	//json.Unmarshal(b, &i)
	links, ok := i["_links"].(map[string]interface{})
	if !ok {
		return "", errors.New("Error converting links")
	}
	next, ok := links["next"].(map[string]interface{})
	if !ok {
		return "", errors.New("Error converting next")
	}
	hrefNext, ok := next["href"].(string)
	if !ok {
		return "", errors.New("Error converting href next")
	}

	return hrefNext, nil
}
func GetSelf(i map[string]interface{}) (string, error) {
	//i := make(map[string]interface{})
	//json.Unmarshal(b, &i)
	links, ok := i["_links"].(map[string]interface{})
	if !ok {
		return "", errors.New("Error converting links")
	}
	self, ok := links["self"].(map[string]interface{})
	if !ok {
		return "", errors.New("Error converting next")
	}
	hrefSelf, ok := self["href"].(string)
	if !ok {
		return "", errors.New("Error converting href next")
	}

	return hrefSelf, nil
}
func GetPrev(i map[string]interface{}) (string, error) {
	//i := make(map[string]interface{})
	//json.Unmarshal(b, &i)
	links, ok := i["_links"].(map[string]interface{})
	if !ok {
		return "", errors.New("Error converting links")
	}
	prev, ok := links["prev"].(map[string]interface{})
	if !ok {
		return "", errors.New("Error converting next")
	}
	hrefPrev, ok := prev["href"].(string)
	if !ok {
		return "", errors.New("Error converting href next")
	}

	return hrefPrev, nil
}
func CountFreeze(i map[string]interface{}, start int64) (*big.Int, int, error) {
	emb, ok := i["_embedded"].(map[string]interface{})
	if !ok {
		return nil, 0, errors.New("cant convert embedded")
	}
	if emb["records"] == nil {
		fmt.Println("-->void")
		return big.NewInt(-1), 0, nil
	}
	recs := emb["records"].([]interface{})
	numbrerOf := len(recs)
	totHere := big.NewInt(0)
	for i := (int)(start); i < len(recs); i++ {
		rec := recs[i].(map[string]interface{})
		state, val, err := CountBosFreeze(rec)
		if err != nil {
			return nil, 0, err
		}
		if state == frozen || state == melting {
			valHere := big.NewInt(val)
			totHere.Add(totHere, valHere)
		}
	}
	fmt.Println("-->")
	//spew.Dump(totHere)
	//spew.Dump(numbrerOf)
	return totHere, numbrerOf, nil
}
func CountBosFreeze(i map[string]interface{}) (FrozenState, int64, error) {
	state, ok := i["state"].(string)
	if !ok {
		return frozen, 0, errors.New("cant convrt state")
	}
	amount, ok := i["amount"].(string)
	if !ok {
		return frozen, 0, errors.New("cant convrt amount")
	}
	amInt, err := strconv.ParseInt(amount, 10, 64)
	//fmt.Println("res", (amInt % 1000000))
	if amInt%1000000 < 500000 {
		amInt = amInt / 1000000
	} else {
		amInt = amInt/1000000 + 1

	}
	if err != nil {
		return frozen, 0, err
	}
	//spew.Dump(amInt)
	switch FrozenState(state) {
	case frozen:
		{
			return frozen, amInt, nil
		}
	case melting:
		{
			return melting, amInt, nil
		}
	case returned:
		{
			return returned, amInt, nil
		}
	default:
		{
			return frozen, 0, nil
		}
	}
}
func RetriveAndCalc(data map[string]interface{}) map[string]interface{} {
	/*data := map[string]interface{}{
		"total":  0,
		"txAddr": "",
		"txNum":  0,
	}*/
	tmp := big.NewInt(0)
	var f map[string]interface{}
	var tot *big.Int
	//ten6 := big.NewInt(1000000)
	var err error
	var hrefNext string
	if data["txAddr"] == "" {
		tot = big.NewInt(0)
		f, err = GetFrozenInter("/api/v1/frozen-accounts")
	} else {
		tot = big.NewInt(data["total"].(int64))
		//n.Mul(n, ten6)

		fmt.Println(tot)

		f, err = GetFrozenInter(data["txAddr"].(string))
		hrefNext = data["txAddr"].(string)
	}
	if err != nil {
		fmt.Println(err)
	}
	/*var out bytes.Buffer
	json.Indent(&out, b, " ", "    ")
	SaveToFile("frozen.json", out.Bytes())*/
	var numberOf int

	tmp, numberOf, err = CountFreeze(f, data["txNum"].(int64))
	if err != nil {
		fmt.Println(err)
	}

	if tmp.Int64() != -1 {
		tot.Add(tot, tmp)
	}
	for tmp.Sign() >= 0 {
		//for i := 0; i < 5; i++ {
		hrefNext, err = GetNext(f)
		if err != nil {
			fmt.Println(err)
		}
		//hrefPrev := GetPrev(f)
		f, err = GetFrozenInter(hrefNext)
		if err != nil {
			fmt.Println(err)
		}
		tmp, numberOf, err = CountFreeze(f, 0)
		if err != nil {
			fmt.Println(err)
		}

		if tmp.Int64() != -1 {
			tot.Add(tot, tmp)
		}
	}
	fmt.Println("toot:", tot)
	fmt.Println("num:", numberOf)
	fmt.Println("ref:", hrefNext)
	intExp := tot.Int64()
	ret := map[string]interface{}{
		"total":  intExp,
		"txAddr": hrefNext,
		"txNum":  numberOf,
	}
	return ret
}
