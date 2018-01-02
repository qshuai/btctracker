package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/tidwall/gjson"
	"go4.org/sort"
)

type BitcoinSite struct {
	Site string
	T    string
}

type TxInfo struct {
	IsIN     bool
	Updated  bool
	Type     string
	TxPrefix string
	AdPrefix string
	TxID     string
	Date     string
	Amount   float64
}

var siteInfo []BitcoinSite = []BitcoinSite{
	{"12Hr4LUeUviQTrHsJtJgy8BgCQ7kk9Nvz2", "BCH"},
	{"1JSX39z6ZuKhUzuAS5QYeLAqstQYXuYAzC", "BCH"},
	{"171C8gR2T5NEU8kMYTo3vTAioDUAhr5ACe", "BCH"},
	{"19tpkUX8MYsSzoP7RK8fL8QDBuZrnXterF", "BCH"},
	{"18ceQ5wyKyUzvycyEzuc7bn7hbc6G1M6iT", "BCH"},
	{"19R7YqtA2hTW8dbswBtbDSai8HMd9hUdEE", "BCH"},
	{"1PFjk6W1R5eGQGyHBgDLJ8wC653a59ZuXo", "BCH"},
	{"1MTLZMBUHDM3mtgGswhcDDCZ1afLJ7FtkE", "BTC"},
	{"1caSJZrCqzRtXgqaSzTkAuMhZK1Cm9tqy", "BTC"},
	{"112xgV6Ejyn5tq8jRryJD28VZmBvnNnYY7", "BTC"},
	{"12A9icBypaLsMBZdZEPpMnRYQFU66GYPvo", "BTC"},
	{"1A9ayyRiWLRwYN7TCVRGNAKvnrQxV8eKRE", "BTC"},
	{"114DFBtAXFrL4S4YxWfXEPyBQpMAQJMvUC", "BTC"},
	{"1KrCGoiPUR4Q2yH3aZcWq1ZXZUieY6D3Qx", "BTC"},
	{"1FzwvTsVtEmyndyrcYFM3WrjotaUezkRW8", "BTC"},
	{"1PzUcJUGSX6D3vaejWMqeryAXThrWGAQCx", "BTC"},
	{"1EMBJyFzA5WW1NBfGHDt8RKaGdZhMbBjqy", "BTC"},
	{"1MtQqfoBQTH28aREpSwmxYo1dFkqGW22t6", "BTC"},
	{"1NCaPoE5UTZpnLq2UXZ3296Yf4YRvnbqZM", "BTC"},
	{"17hxK6AwSVBVqMheYNm8qmdzvS65abXmPZ", "BTC"},
	{"1C6n9guftW3dZ3EHHsMGHvETLqnSTeANn5", "BTC"},
	{"1HetHPqXPcEKv1eePp1FLusLcmxA2a1hez", "BTC"},
	{"14qp5GhyezEz3N3NDceQW6EHqffcBXw9d3", "BTC"},
	{"191yAE7Mwz4nGLgTydigbKn2AGNEPXGkYe", "BTC"},
	{"1FYr29v8eVQ26WUruWA2KC9zuVShG3MbMd", "BTC"},
	{"1DoAGgm3zAcPUtJje9ZyM9jZyAkL26nYA6", "BTC"},
	{"1G2QbTNhYDhiQ8s6RJJMBD4sE8JffV5Pso", "BTC"},
	{"1FhdGQ8tCFkgfeFqyGHrYNJEJXMQ2bVGjy", "BTC"},
}

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	var info = make(map[string][]TxInfo)

	length := len(siteInfo)
	for i := 0; i < length; i++ {
		value := siteInfo[i]
		var prefix string
		if value.T == "BTC" {
			prefix = "https://chain-vip-bj.api.btc.com/v3/address/"
		} else {
			prefix = "https://bcc-chain-vip-bj.api.btc.com/v3/address/"
		}

		content := getContet(prefix + value.Site)

		hash := gjson.Get(content, "data.list.#.hash").Array()
		amount := gjson.Get(content, "data.list.#.inputs_value").Array()
		created := gjson.Get(content, "data.list.#.created_at").Array()
		inputs := gjson.Get(content, "data.list.#.inputs").Array()

		for index, iterm := range hash {
			txinfo := TxInfo{}
			if value.T == "BTC" {
				txinfo.TxPrefix = "https://btc.com/"
				txinfo.AdPrefix = "https://btc.com/"
				txinfo.Type = "BTC"
			} else {
				txinfo.TxPrefix = "https://bch.btc.com/"
				txinfo.AdPrefix = "https://bch.btc.com/"
				txinfo.Type = "BCH"
			}
			txinfo.Amount = amount[index].Float() / float64(100000000)
			tm := time.Unix(created[index].Int(), 0)
			txinfo.Date = tm.Format("2006-01-02 15:04:05")
			txinfo.TxID = iterm.String()
			if created[index].Int() > 1514894400 {
				txinfo.Updated = true
			}

			if(!strings.Contains(inputs[index].String(), value.Site)){
				txinfo.IsIN = true
			}

			info[value.Site] = append(info[value.Site], txinfo)
		}
	}

	address := make([]string, 0)
	for key, iterm := range info {
		address = append(address, iterm[0].Date+key)
	}

	sort.Strings(address)
	l := len(address)
	str := make([]string, 0)
	for i := l - 1; i >= 0; i-- {
		str = append(str, address[i][19:])
	}
	c.Data["str"] = str
	c.Data["lists"] = info
	c.TplName = "index.html"
}

func getContet(url string) string {
	get, err := http.Get(url + "/tx")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer get.Body.Close()
	content, _ := ioutil.ReadAll(get.Body)
	return string(content)
}

/**
[
	[
		{"addresses":["1KrNNBKZ1eovW7anXHhfv95drYwMKE7ayx"],
		"value":129143420,
		"type":"P2PKH",
		"spent_by_tx":"4af0cea3ec7af8b4843726080d39b356805ef2cd70918ee25d50956b97dff3ef",
		"spent_by_tx_position":53
		},
		{"addresses":["1JrsYGFSqvEArDsNAt96ViTTBsuJvypQwp"],
		"value":3511815315,
		"type":"P2PKH",
		"spent_by_tx":"e9be116ac1f0894b822797ce989a5c05c4cdb3245ecc2fbb39fdf4582209ba60",
		"spent_by_tx_position":0
		}
	]
	[
		{"addresses":["171C8gR2T5NEU8kMYTo3vTAioDUAhr5ACe"],
		"value":3640962351,
		"type":"P2PKH",
		"spent_by_tx":"c3bb0ee836b0f2eac6b23c7ad12b9a35a50fe36b97a7c3c167ef6255ff3e62eb",
		"spent_by_tx_position":0
		}
	]
]
*/
