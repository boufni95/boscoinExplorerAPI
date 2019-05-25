package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type FrozenState string

const (
	frozen   FrozenState = "frozen"
	melting  FrozenState = "melting"
	returned FrozenState = "returned"
)

func GetFrozenBytes(moreP string) []byte {

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
	return received
}
func GetFrozenInter(moreP string) map[string]interface{} {

	received := GetFrozenBytes(moreP)
	i := make(map[string]interface{})
	json.Unmarshal(received, &i)
	return i
}

func GetNext(i map[string]interface{}) string {
	//i := make(map[string]interface{})
	//json.Unmarshal(b, &i)
	links := i["_links"].(map[string]interface{})
	next := links["next"].(map[string]interface{})
	hrefNext := next["href"].(string)

	return hrefNext
}
func GetSelf(i map[string]interface{}) string {
	//i := make(map[string]interface{})
	//json.Unmarshal(b, &i)
	links := i["_links"].(map[string]interface{})
	self := links["self"].(map[string]interface{})
	hrefNext := self["href"].(string)

	return hrefNext
}
func GetPrev(i map[string]interface{}) string {
	//i := make(map[string]interface{})
	//json.Unmarshal(b, &i)
	links := i["_links"].(map[string]interface{})
	prev := links["prev"].(map[string]interface{})
	hrefNext := prev["href"].(string)

	return hrefNext
}
func CountFreeze(i map[string]interface{}) (*big.Int, int) {
	emb := i["_embedded"].(map[string]interface{})
	if emb["records"] == nil {
		fmt.Println("-->void")
		return big.NewInt(-1), 0
	}
	recs := emb["records"].([]interface{})
	numbrerOf := len(recs)
	totHere := big.NewInt(0)
	for i, _ := range recs {
		rec := recs[i].(map[string]interface{})
		state, val := CountBosFreeze(rec)
		if state == frozen || state == melting {
			valHere := big.NewInt(val)
			totHere.Add(totHere, valHere)
		}
	}
	fmt.Println("-->")
	spew.Dump(totHere)
	spew.Dump(numbrerOf)
	return totHere, numbrerOf
}
func CountBosFreeze(i map[string]interface{}) (FrozenState, int64) {
	state := i["state"].(string)
	amount := i["amount"].(string)
	amInt, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return frozen, 0
	}
	//spew.Dump(amInt)
	switch FrozenState(state) {
	case frozen:
		{
			return frozen, amInt
		}
	case melting:
		{
			return melting, amInt
		}
	case returned:
		{
			return returned, amInt
		}
	default:
		{
			return frozen, 0
		}
	}
}
