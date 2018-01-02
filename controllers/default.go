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

//site information
type BitcoinSite struct {
	Site string
	T    string
}

//tx relational information
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

//read configuration
var sites []string = beego.AppConfig.Strings("watchsite")
var types []string = beego.AppConfig.Strings("watchtype")

//define global variable
var siteInfo []BitcoinSite

//just do this
type MainController struct {
	beego.Controller
}

//Get Method for data list
func (c *MainController) Index() {
	//create empty map for store data
	var info = make(map[string][]TxInfo)

	//get all sites
	for index,value := range sites{
		siteInfo = append(siteInfo, BitcoinSite{value, types[index]})
	}

	//
	length := len(siteInfo)
	for i := 0; i < length; i++ {
		value := siteInfo[i]
		var prefix string
		if value.T == "BTC" {
			prefix = beego.AppConfig.String("btcprefix")
		} else {
			prefix = beego.AppConfig.String("bchprfix")
		}

		content := getContet(prefix + value.Site)
		time.Sleep(500*time.Millisecond)

		hash := gjson.Get(content, "data.list.#.hash").Array()

		//if you just want inputs_value, please justify as following:
		//amount := gjson.Get(content, "data.list.#.inputs_value").Array()
		amount := gjson.Get(content, "data.list.#.outputs_value").Array()
		created := gjson.Get(content, "data.list.#.created_at").Array()
		inputs := gjson.Get(content, "data.list.#.inputs").Array()

		//assembly data
		for index, iterm := range hash {
			txinfo := TxInfo{}
			//just support BTC and BCH
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

			//justify whether income or expense
			if(!strings.Contains(inputs[index].String(), value.Site)){
				txinfo.IsIN = true
			}

			//get final result
			info[value.Site] = append(info[value.Site], txinfo)
		}

	}

	//sort map indirectly
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

	//delivery to template
	c.Data["str"] = str
	c.Data["lists"] = info
	c.TplName = "index.html"
}

//get http raw data
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
