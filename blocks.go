package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetBlocksBytes(cursor uint64, limit uint64) []byte {
	scursor := strconv.FormatUint(cursor, 10)
	slimit := strconv.FormatUint(limit, 10)

	var path strings.Builder
	path.WriteString("https://mainnet.blockchainos.org/api/v1/blocks?cursor=")
	path.WriteString(scursor)
	path.WriteString("&limit=")
	path.WriteString(slimit)
	path.WriteString("&reverse=false")

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
func GetBlockBytes(hash string) []byte {

	var path strings.Builder
	path.WriteString("https://mainnet.blockchainos.org/api/v1/blocks/")
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
