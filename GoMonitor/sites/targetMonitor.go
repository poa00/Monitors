package TargetMonitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Autogenerated struct {
	Data struct {
		Product struct {
			Typename        string `json:"__typename"`
			Tcin            string `json:"tcin"`
			NotifyMeEnabled bool   `json:"notify_me_enabled"`
			StorePositions  []struct {
				Aisle int    `json:"aisle"`
				Block string `json:"block"`
			} `json:"store_positions"`
			Fulfillment struct {
				ProductID                       string `json:"product_id"`
				IsOutOfStockInAllStoreLocations bool   `json:"is_out_of_stock_in_all_store_locations"`
				StoreOptions                    []struct {
					LocationName                       string  `json:"location_name"`
					LocationAddress                    string  `json:"location_address"`
					LocationID                         string  `json:"location_id"`
					SearchResponseStoreType            string  `json:"search_response_store_type"`
					LocationAvailableToPromiseQuantity float64 `json:"location_available_to_promise_quantity"`
					OrderPickup                        struct {
						PickupDate         string `json:"pickup_date"`
						GuestPickSLA       int    `json:"guest_pick_sla"`
						AvailabilityStatus string `json:"availability_status"`
					} `json:"order_pickup"`
					InStoreOnly struct {
						AvailabilityStatus string `json:"availability_status"`
					} `json:"in_store_only"`
				} `json:"store_options"`
				ShippingOptions struct {
					AvailabilityStatus         string  `json:"availability_status"`
					LoyaltyAvailabilityStatus  string  `json:"loyalty_availability_status"`
					AvailableToPromiseQuantity float64 `json:"available_to_promise_quantity"`
					MinimumOrderQuantity       float64 `json:"minimum_order_quantity"`
					Services                   []struct {
						IsTwoDayShipping               bool      `json:"is_two_day_shipping"`
						IsBaseShippingMethod           bool      `json:"is_base_shipping_method"`
						ShippingMethodID               string    `json:"shipping_method_id"`
						Cutoff                         time.Time `json:"cutoff"`
						MaxDeliveryDate                string    `json:"max_delivery_date"`
						MinDeliveryDate                string    `json:"min_delivery_date"`
						ShippingMethodShortDescription string    `json:"shipping_method_short_description"`
						ServiceLevelDescription        string    `json:"service_level_description"`
					} `json:"services"`
				} `json:"shipping_options"`
				ScheduledDelivery struct {
					AvailabilityStatus string `json:"availability_status"`
				} `json:"scheduled_delivery"`
			} `json:"fulfillment"`
		} `json:"product"`
	} `json:"data"`
}
type MyJsonName struct {
	Data struct {
		Product struct {
			Typename string `json:"__typename"`
			Item     struct {
				Enrichment struct {
					BuyURL string `json:"buy_url"`
					Images struct {
						AlternateImageUrls []string `json:"alternate_image_urls"`
						AlternateImages    []string `json:"alternate_images"`
						BaseURL            string   `json:"base_url"`
						ContentLabels      []struct {
							ImageURL string `json:"image_url"`
						} `json:"content_labels"`
						PrimaryImage    string `json:"primary_image"`
						PrimaryImageURL string `json:"primary_image_url"`
					} `json:"images"`
				} `json:"enrichment"`
				ProductDescription struct {
					BulletDescriptions    []string `json:"bullet_descriptions"`
					DownstreamDescription string   `json:"downstream_description"`
					SoftBulletDescription string   `json:"soft_bullet_description"`
					SoftBullets           struct {
						Bullets []string `json:"bullets"`
						Title   string   `json:"title"`
					} `json:"soft_bullets"`
					Title string `json:"title"`
				} `json:"product_description"`
				ProductVendors []struct {
					ID         string `json:"id"`
					VendorName string `json:"vendor_name"`
				} `json:"product_vendors"`
				RelationshipTypeCode       string   `json:"relationship_type_code"`
				ReturnPoliciesGuestMessage string   `json:"return_policies_guest_message"`
				Ribbons                    []string `json:"ribbons"`
			} `json:"item"`
			Price struct {
				CurrentRetail             float64 `json:"current_retail"`
				FormattedCurrentPrice     string  `json:"formatted_current_price"`
				FormattedCurrentPriceType string  `json:"formatted_current_price_type"`
				IsCurrentPriceRange       bool    `json:"is_current_price_range"`
			} `json:"price"`
			TargetFinds struct{} `json:"target_finds"`
			Tcin        string   `json:"tcin"`
		} `json:"product"`
	} `json:"data"`
}
type Config struct {
	sku              string
	skuName          string // Only for new Egg
	startDelay       int
	discord          string
	site             string
	priceRangeMax    int
	priceRangeMin    int
	proxyCount       int
	indexMonitorJson int
}
type Monitor struct {
	Config              Config
	monitorProduct      Product
	Availability        bool
	currentAvailability bool
	Client              http.Client
	file                *os.File
}
type Product struct {
	name        string
	stockNumber int
	offerId     string
	price       int
	image       string
}
type Proxy struct {
	ip       string
	port     string
	userAuth string
	userPass string
}
type ItemInMonitorJson struct {
	Sku  string `json:"sku"`
	Site string `json:"site"`
	Stop bool   `json:"stop"`
	Name string `json:"name"`
}


var file os.File

// func walmartMonitor(sku string) {
// 	go NewMonitor(sku, 1, 1000)
// 	fmt.Scanln()
// }

func NewMonitor(sku string, priceRangeMin int, priceRangeMax int) *Monitor {
	defer func() {
	     if r := recover(); r != nil {
	        fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
	    }
	  }()
	fmt.Println("TESTING", sku, priceRangeMin, priceRangeMax)
	m := Monitor{}
	m.Availability = false
	var err error
	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.site = "Target"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 60 * time.Second}
	m.Config.discord = "https://discord.com/api/webhooks/826281048480153641/rmifnt8w6NKFainUqAsE16RZM1LzNGrPdUB0jP5M3PJwm0hRvRmemyrqr0FdrZEBMOmd"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = 10
	m.Config.priceRangeMax = priceRangeMax
	m.Config.priceRangeMin = priceRangeMin
	m.getProductImage(sku)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	defer file.Close()

	data, err := ioutil.ReadFile("GoMonitors.json")
	if err != nil {
		fmt.Print(err)
	}
	// fmt.Println(string(data))
	var monitorCheckJson []interface{}
	err = json.Unmarshal(data, &monitorCheckJson)
	fmt.Println(monitorCheckJson)
	for key, value := range monitorCheckJson {
		var currentObject ItemInMonitorJson
		currentObject.Site = value.(map[string]interface{})["site"].(string)
		currentObject.Stop = value.(map[string]interface{})["stop"].(bool)
		currentObject.Name = value.(map[string]interface{})["name"].(string)
		currentObject.Sku = value.(map[string]interface{})["sku"].(string)
		if currentObject.Sku == m.Config.sku {
			m.Config.indexMonitorJson = key
			fmt.Println(currentObject, key)
		}
	}

	path := "cloud.txt"
	var proxyList = make([]string, 0)
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	snl := bufio.NewScanner(buf)
	for snl.Scan() {
		proxy := snl.Text()
		proxyList = append(proxyList, proxy)
		splitProxy := strings.Split(string(proxy), ":")
		newProxy := Proxy{}
		newProxy.userAuth = splitProxy[2]
		newProxy.userPass = splitProxy[3]
		newProxy.ip = splitProxy[0]
		newProxy.port = splitProxy[1]
		//	go NewMonitor(newProxy)
		//	time.Sleep(5 * time.Second)
	}
	err = snl.Err()
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	i := true
	for i == true {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
			}
		}()
		data, err := ioutil.ReadFile("GoMonitors.json")
		if err != nil {
			fmt.Print(err)
		}
		// fmt.Println(string(data))
		var monitorCheckJson []interface{}
		err = json.Unmarshal(data, &monitorCheckJson)
		var currentObject ItemInMonitorJson
		currentObject.Site = monitorCheckJson[m.Config.indexMonitorJson].(map[string]interface{})["site"].(string)
		currentObject.Stop = monitorCheckJson[m.Config.indexMonitorJson].(map[string]interface{})["stop"].(bool)
		currentObject.Name = monitorCheckJson[m.Config.indexMonitorJson].(map[string]interface{})["name"].(string)
		currentObject.Sku = monitorCheckJson[m.Config.indexMonitorJson].(map[string]interface{})["sku"].(string)
		if !currentObject.Stop {
			currentProxy := m.getProxy(proxyList)
			splittedProxy := strings.Split(currentProxy, ":")
			proxy := Proxy{splittedProxy[0], splittedProxy[1], splittedProxy[2], splittedProxy[3]}
			//	fmt.Println(proxy, proxy.ip)
			prox1y := fmt.Sprintf("http://%s:%s@%s:%s", proxy.userAuth, proxy.userPass, proxy.ip, proxy.port)
			proxyUrl, err := url.Parse(prox1y)
			if err != nil {
				fmt.Println(err)
				m.file.WriteString(err.Error() + "\n")
				return nil
			}
			defaultTransport := &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
			m.Client.Transport = defaultTransport
			go m.monitor()
			time.Sleep(500 * (time.Millisecond))
			fmt.Println(m.Availability)
		} else {
			fmt.Println(currentObject.Sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *Monitor) monitor() error {
	fmt.Println("Monitoring")
		defer func() {
	     if r := recover(); r != nil {
	        fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
	    }
	  }()
	// url := "https://httpbin.org/ip"

	// req, _ := http.NewRequest("GET", url, nil)

	// res, _ := m.Client.Do(req)

	// defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	// fmt.Println(string(body))

	url := fmt.Sprintf("https://redsky.target.com/redsky_aggregations/v1/web/pdp_fulfillment_v1?key=ff457966e64d5e877fdbad070f276d18ecec4a01&tcin=%s&store_id=2067&store_positions_store_id=2067&has_store_positions_store_id=true&scheduled_delivery_store_id=2067&pricing_store_id=2067&fulfillment_test_mode=grocery_opu_team_member_test", m.Config.sku)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	// req.Header.Add("cookie", "TealeafAkaSid=r5S-XRsuxWbk94tkqVB3CruTmaJKz32Z")
	res, err := m.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	//	fmt.Println(res)
	//	fmt.Println(string(body))
	fmt.Println(res.StatusCode)
	var realBody Autogenerated
	err = json.Unmarshal(body, &realBody)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}

	currentAvailability := realBody.Data.Product.Fulfillment.ShippingOptions.AvailabilityStatus
	if currentAvailability == "PRE_ORDER_UNSELLABLE" || currentAvailability == "UNAVAILABLE" || currentAvailability == "OUT_OF_STOCK" {
                        m.currentAvailability = false
    } else if currentAvailability == "PRE_ORDER_SELLABLE" || currentAvailability == "IN_STOCK" {
                        m.currentAvailability = true
	} else {
		m.currentAvailability = false
		fmt.Println(currentAvailability)
	}
	m.monitorProduct.stockNumber = int(realBody.Data.Product.Fulfillment.ShippingOptions.AvailableToPromiseQuantity)
//	m.monitorProduct.price = realBody.Data.Product.Fulfillment.

	// log.Printf("%+v", m.Availability)
	if m.Availability == false && m.currentAvailability == true {
		fmt.Println("Item in Stock")
		m.sendWebhook()

	}
	if m.Availability == true && m.currentAvailability == false {
		fmt.Println("Item Out Of Stock")
	}
	m.Availability = m.currentAvailability
	return nil
}

func (m *Monitor) getProxy(proxyList []string) string {
defer func() {
	     if r := recover(); r != nil {
	        fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
	    }
	  }()
	//fmt.Scanln()
	// rand.Seed(time.Now().UnixNano())
	// randomPosition := rand.Intn(len(proxyList)-0) + 0
	if m.Config.proxyCount+1 == len(proxyList) {
		m.Config.proxyCount = 0
	}
	m.Config.proxyCount++
	//fmt.Println(proxyList[m.Config.proxyCount])
	return proxyList[m.Config.proxyCount]
}

func (m *Monitor) sendWebhook() error {
	defer func() {
	     if r := recover(); r != nil {
	        fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
	    }
	  }()
	for _, letter := range m.monitorProduct.name {
		if string(letter) == `"` {
			m.monitorProduct.name = strings.Replace(m.monitorProduct.name, `"`, "", -1)
		}
	}
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	payload := strings.NewReader(fmt.Sprintf(`{
  "content": null,
  "embeds": [
    {
      "title": "%s Monitor",
      "url": "https://www.target.com/p/prada/-/A-%s",
      "color": 507758,
      "fields": [
        {
          "name": "Product Name",
          "value": "%s"
        },
        {
          "name": "Product Availability",
          "value": "In Stock",
          "inline": true
        },
        {
          "name": "Price",
          "value": "%d",
          "inline": true
        },
        {
          "name": "Tcin",
          "value": "%s",
          "inline": true
        },
        {
          "name": "Stock Number",
          "value": "%d"
        },
        {
          "name": "Links",
          "value": "[Product](https://www.target.com/p/prada/-/A-%ss) | [Cart](https://www.target.com/co-cart)"
        }
      ],
      "footer": {
        "text": "Prada#4873"
      },
      "timestamp": "2021-04-01T18:40:00.000Z",
      "thumbnail": {
        "url": "%s"
      }
    }
  ],
  "avatar_url": "https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png"
}`, m.Config.site, m.Config.sku, m.monitorProduct.name, m.monitorProduct.price, m.Config.sku, m.monitorProduct.stockNumber, m.Config.sku, m.monitorProduct.image))
	req, err := http.NewRequest("POST", m.Config.discord, payload)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("accept", "application/json")
	req.Header.Add("dnt", "1")
	req.Header.Add("accept-language", "en")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("sec-fetch-site", "cross-site")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	fmt.Println(res)
	fmt.Println(string(body))
	fmt.Println(payload)
	return nil
}

func (m *Monitor) getProductImage (tcin string) {
	defer func() {
	     if r := recover(); r != nil {
	        fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
	    }
	  }()
	fmt.Println("Getting Product Image")
	url := fmt.Sprintf("https://redsky.target.com/redsky_aggregations/v1/web/pdp_client_v1?key=ff457966e64d5e877fdbad070f276d18ecec4a01&tcin=%s&member_id=20032056835&store_id=2067&has_store_id=true&pricing_store_id=2067&has_pricing_store_id=true&scheduled_delivery_store_id=2067&has_scheduled_delivery_store_id=true&has_financing_options=false&visitor_id=0178E7894D540201A352D90ED642CB06&has_size_context=true", tcin)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return 
	}
	res, err := m.Client.Do(req)
		if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return 
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
		if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return 
	}
	fmt.Println(res.StatusCode)
	var realBody MyJsonName
	err = json.Unmarshal(body, &realBody)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return
	}
	m.monitorProduct.name = realBody.Data.Product.Item.ProductDescription.Title
	m.monitorProduct.image = realBody.Data.Product.Item.Enrichment.Images.PrimaryImageURL
	m.monitorProduct.price = int(realBody.Data.Product.Price.CurrentRetail)
	fmt.Println(m.monitorProduct.image)
	return
}