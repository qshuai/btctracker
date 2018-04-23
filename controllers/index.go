package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tidwall/gjson"
	"go4.org/sort"
)

// site information
type BitcoinSite struct {
	Site string
	T    string
}

// tx relational information
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

// read configuration
var sites = beego.AppConfig.Strings("watchsite")
var types = beego.AppConfig.Strings("watchtype")
var hasUpdated = make(map[string]struct{})

var self *MainController
var once sync.Once

// just do this
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

func StoreAllBitSite() {
	// get all sites
	for index, value := range sites {
		self.siteInfo = append(self.siteInfo, BitcoinSite{value, types[index]})
	}
}

// Get Method for data list
func (c *MainController) Index() {
	controller := GetMainController()
	controller.lock.RLock()
	defer controller.lock.RUnlock()

	// delivery to template
	c.Data["str"] = controller.str
	c.Data["lists"] = controller.info
	c.TplName = "index.html"
}

// get http raw data
func getContent(url string) string {
	get, err := http.Get(url + "/tx")
	if err != nil {
		logs.Error(err)
		return ""
	}
	defer get.Body.Close()
	content, _ := ioutil.ReadAll(get.Body)
	return string(content)
}

// store data
func (c *MainController) StoreDate() {
	info := make(map[string][]TxInfo)
	str := make([]string, 0)
	isUpdated := false

	length := len(c.siteInfo)
	for i := 0; i < length; i++ {
		time.Sleep(1 * time.Second)
		value := c.siteInfo[i]
		var prefix string
		if value.T == "BTC" {
			prefix = beego.AppConfig.String("btcprefix")
		} else {
			prefix = beego.AppConfig.String("bchprefix")
		}

		content := getContent(prefix + value.Site)

		hash := gjson.Get(content, "data.list.#.hash").Array()

		// if you just want inputs_value, please justify as following:
		// amount := gjson.Get(content, "data.list.#.inputs_value").Array()
		amount := gjson.Get(content, "data.list.#.outputs_value").Array()
		created := gjson.Get(content, "data.list.#.created_at").Array()
		inputs := gjson.Get(content, "data.list.#.inputs").Array()

		// assembly data
		for index, item := range hash {
			txinfo := TxInfo{}
			// just support BTC and BCH
			if value.T == "BTC" {
				txinfo.TxPrefix = "https://btc.com/"
				txinfo.AdPrefix = "https://btc.com/"
				txinfo.Type = "BTC"
			} else {
				txinfo.TxPrefix = "https://bch.btc.com/"
				txinfo.AdPrefix = "https://bch.btc.com/"
				txinfo.Type = "BCH"
			}
			txinfo.Amount = amount[index].Float() / float64(1e8)
			tm := time.Unix(created[index].Int(), 0)
			txinfo.Date = tm.Format("2006-01-02 15:04:05")
			txinfo.TxID = item.String()
			if created[index].Int() > 1514894400 {
				if _, ok := hasUpdated[item.String()]; !ok {
					logs.Alert("Alert: address's info has updated(%s:%s)", value.T, item)
					hasUpdated[item.String()] = struct{}{}
				}

				isUpdated = true
				txinfo.Updated = true
			}

			// justify whether income or expense
			if !strings.Contains(inputs[index].String(), value.Site) {
				txinfo.IsIN = true
			}

			// get final result
			info[value.Site] = append(info[value.Site], txinfo)
		}

	}

	// sort map indirectly
	address := make([]string, 0)
	for key, item := range info {
		address = append(address, item[0].Date+key)
	}
	sort.Strings(address)
	l := len(address)

	for i := l - 1; i >= 0; i-- {
		str = append(str, address[i][19:])
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.info = info
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
