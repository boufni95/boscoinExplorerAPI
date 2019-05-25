package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func GetTxBytes(hash string) []byte {

	var path strings.Builder
	path.WriteString("https://mainnet.blockchainos.org/api/v1/transactions/")
	path.WriteString(hash)

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
func ExtrTxsBlocks(i map[string]interface{}) (int, int, error) {

	emb, ok := i["_embedded"].(map[string]interface{})
	if !ok {
		return 0, 0, errors.New("error on embedded")
	}
	recs, ok := emb["records"].([]interface{})
	if !ok {
		return 0, 0, errors.New("error on records")
	}
	l := len(recs)
	cont := 0
	for _, v := range recs {
		myv, ok := v.(map[string]interface{})
		if !ok {
			return 0, 0, errors.New("error on txs")
		}
		if myv["transactions"] != nil {
			txs := myv["transactions"].([]interface{})
			for _, v := range txs {
				spew.Dump(v)
				b := GetTxBytes((v).(string))
				pri := FmtIndent(b)
				fmt.Println(string(pri))
			}
			cont++

		}
		//spew.Dump(myv["transactions"])
		//fmt.Println("_______________________________")
	}
	return l, cont, nil
}
