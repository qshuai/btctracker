package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
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

var self *MainController
var once sync.Once

//just do this
type MainController struct {
	beego.Controller

	siteInfo  []BitcoinSite
	info      map[string][]TxInfo
	str       []string
	IsUpdated bool
	lock      sync.RWMutex
}

func GetMainController() *MainController {
	once.Do(func() {
		self = &MainController{
			siteInfo:  []BitcoinSite{},
			info:      make(map[string][]TxInfo),
			str:       []string{},
			IsUpdated: false,
			lock:      sync.RWMutex{},
		}
	})

	return self
}

//Get Method for data list
func (c *MainController) Index() {
	controller := GetMainController()
	controller.lock.RLock()
	defer controller.lock.RUnlock()

	//delivery to template
	c.Data["str"] = controller.str
	c.Data["lists"] = controller.info
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

//store data
func (c *MainController) StoreDate() {
	siteInfo := []BitcoinSite{}
	info := make(map[string][]TxInfo)
	str := []string{}
	isUpdated := false

	//get all sites
	for index, value := range sites {
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

		//API request rate limit
		time.Sleep(500 * time.Millisecond)

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
				isUpdated = true
				txinfo.Updated = true
			}

			//justify whether income or expense
			if !strings.Contains(inputs[index].String(), value.Site) {
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

	for i := l - 1; i >= 0; i-- {
		str = append(str, address[i][19:])
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.info = info
	c.siteInfo = siteInfo
	c.str = str
	c.IsUpdated = isUpdated
}

func (c *MainController) Timer(duration time.Duration) {
	t := time.NewTimer(duration)
	for {
		select {
		case <-t.C:
			c.StoreDate()
			t.Reset(duration)
		}
	}
}
