package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func GetBytes(moreP string) []byte {

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

func SaveToFile(fileName string, b []byte) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(fileName, b, 0644)
}
func FmtIndent(b []byte) []byte {
	var out bytes.Buffer
	json.Indent(&out, b, " ", "    ")
	return out.Bytes()
}
