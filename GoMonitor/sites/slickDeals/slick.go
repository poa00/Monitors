package slickDealsMonitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Config struct {
	sku              string
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
	Availability        string
	currentAvailability string
	Client              http.Client
	file                *os.File
	stop                bool
	sku                 []string
	CurrentCompanies    []Company
}
type Product struct {
	name        string
	stockNumber string
	productId   string
	price       int
	image       string
	link        string
}
type Proxy struct {
	ip       string
	port     string
	userAuth string
	userPass string
}
type ItemInMonitorJson struct {
	Sku       string `json:"sku"`
	Site      string `json:"site"`
	Stop      bool   `json:"stop"`
	Name      string `json:"name"`
	Companies []Company
}
type Company struct {
	Company      string `json:"company"`
	Webhook      string `json:"webhook"`
	Color        string `json:"color"`
	CompanyImage string `json:"companyImage"`
}

var file os.File

// func walmartMonitor(sku string) {
// 	go NewMonitor(sku, 1, 1000)
// 	fmt.Scanln()
// }

func NewMonitor() *Monitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : Slick Deals Recovering from panic in printAllOperations error is: %v \n", r)
		}
	}()
	fmt.Println("TESTING")
	m := Monitor{}
	m.Availability = "OUT_OF_STOCK_ONLINE"
	var err error
	//	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.site = "Slick Deals"
	m.Config.startDelay = 3000
	//	m.Config.sku = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.discord = "https://discord.com/api/webhooks/842460357607555075/pKynxp-r1_-VbWAcnhleC2vCKo357vCb29Y83jLG_NgPfiqY3WuUdjZILVWYpQ-SmAj5"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = ""

	path := "cloud.txt"
	var proxyList = make([]string, 0)
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	// defer func() {
	// 	if err = buf.Close(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

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
	buf.Close()
	err = snl.Err()
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	go m.checkStop()
	time.Sleep(5000 * (time.Millisecond))
	i := true
	for i == true {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
			}
		}()
		if !m.stop {
			currentProxy := m.getProxy(proxyList)
			splittedProxy := strings.Split(currentProxy, ":")
			proxy := Proxy{splittedProxy[0], splittedProxy[1], splittedProxy[2], splittedProxy[3]}
			//	fmt.Println(proxy, proxy.ip)
			prox1y := fmt.Sprintf("http://%s:%s@%s:%s", proxy.userAuth, proxy.userPass, proxy.ip, proxy.port)
			proxyUrl, err := url.Parse(prox1y)
			if err != nil {
				fmt.Println(err)

				return nil
			}
			defaultTransport := &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
			m.Client.Transport = defaultTransport
			m.monitor()
			//	time.Sleep(500 * (time.Millisecond))
			fmt.Println("Slick Deals : ", len(m.sku))
		} else {
			fmt.Println(m.Config.sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *Monitor) monitor() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()
	//	fmt.Println("Monitoring")
	// 	defer func() {
	//      if r := recover(); r != nil {
	//         	        fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
	//     }
	//   }()
	// url := "https://httpbin.org/ip"

	// req, _ := http.NewRequest("GET", url, nil)

	// res, _ := m.Client.Do(req)

	// defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	// fmt.Println(string(body))

	url := "https://slickdeals.net/live/spy.php?thread=15023128&post=147235393&threadrate=93848830&time=1620921696&maxitems=20&forum=9"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)

		return nil
	}
	// req.Header.Add("authority", "discord.com")
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
	// req.Header.Add("cookie", "TealeafAkaSid=r5S-XRsuxWbk94tkqVB3CruTmaJKz32Z")
	res, err := m.Client.Do(req)
	if err != nil {
		fmt.Println(err)

		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

		return nil
	}
	//	fmt.Println(res)
	fmt.Println(res.StatusCode)
	if res.StatusCode != 200 {
		return nil
	}

	parsed := string(body)
	threads := strings.Split(parsed, `<![CDATA[<div id="thread_`)
	// fmt.Println(threads[1])
	for _, value := range threads {
		newThread := strings.Split(value, `]]></htmlbit>`)[0]
		id := strings.Split(value, `"`)[0]
		//	newThread := strings.Split(value, "]]>")[0]
		firstDesc := strings.Split(newThread, `<a target="_blank" href="/f/`)
		var link string
		var title string
		var isPresent bool
		for _, v := range m.sku {
			if v == id {
				isPresent = true
			}
		}

		if isPresent != true {
			link = fmt.Sprintf("https://slickdeals.net/f/%s", strings.Split(firstDesc[1], `"`)[0])
			title = strings.Split(strings.Split(firstDesc[1], `" >`)[1], "</a>")[0]
			price := strings.Split(title, "$")
			var returnPrice string

			if len(price) == 2 {
				d := fmt.Sprintf("$%s", returnSplitted(price[1]))
				returnPrice = d
			} else if len(price) == 3 {
				returnPrice = fmt.Sprintf("$%s, $%s", returnSplitted(price[1]), returnSplitted(price[2]))
			} else {
				returnPrice = fmt.Sprintf("$%s, $%s", returnSplitted(price[1]), returnSplitted(price[2]))
			}

			if isPresent != true {
				desc, image, links := m.getDesc(link)
				go m.sendWebhook(link, title, returnPrice, id, desc, image, links)
				m.sku = append(m.sku, id)
			}
		}

		// fmt.Println(len(price))
		// fmt.Println(id,"\n", link, "\n",  title, "\n\n\n\n")

	}

	return nil
}

func returnSplitted(s string) string {
	var returnString string
	returnString = strings.Split(s, " ")[0]
	return returnString
}

func (m *Monitor) getProxy(proxyList []string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
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

func (m *Monitor) sendWebhook(link string, title string, price string, id string, desc string, image string, links string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()
	for _, letter := range m.monitorProduct.name {
		if string(letter) == `"` {
			m.monitorProduct.name = strings.Replace(m.monitorProduct.name, `"`, "", -1)
		}
	}
	fmt.Println("Testing Here : ", m.monitorProduct.name, "Here")
	if strings.HasSuffix(m.monitorProduct.name, "                       ") {
		m.monitorProduct.name = strings.Replace(m.monitorProduct.name, "                       ", "", -1)
	}
	fmt.Println("Testing Here : ", m.monitorProduct.name, "Here")
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	now := time.Now()
	currentTime := strings.Split(now.String(), "-0400")[0]
	if strings.HasSuffix(currentTime, " ") {
		currentTime = strings.TrimSuffix(currentTime, " ")
	}
	re := regexp.MustCompile(`\r?\n`)
	desc = re.ReplaceAllString(desc, `\n`)
	payload := strings.NewReader(fmt.Sprintf(`{
  "content": null,
  "embeds": [
    {
      "title": "%s Monitor",
      "url": "%s",
      "color": 507758,
      "fields": [
        {
          "name": "Product Name",
		  "value": "%s"
		},
		{
			"name": "Product Description",
			"value": "%s"
		  },
        {
          "name": "Price",
          "value": "%s",
          "inline": true
        },
        {
          "name": "Product ID",
          "value": "%s",
          "inline": true
        },
        {
          "name": "Links",
          "value": "[Thread](%s) | [Product](%s)"
        }
      ],
      "footer": {
        "text": "Prada#4873"
      },
      "timestamp": "2021-05-13 13:57:26.5157268",
      "thumbnail": {
        "url": "%s"
      }
    }
  ],
  "avatar_url": "https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png"
}`, m.Config.site, link, title, desc, price, id, link, links, image))
	req, err := http.NewRequest("POST", m.Config.discord, payload)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)

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

		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)

		return nil
	}
	fmt.Println(res)
	fmt.Println(string(body))
	fmt.Println(payload)
	return nil
}

func (m *Monitor) getDesc(link string) (string, string, string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("authority", "slickdeals.net")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("accept", "*/*")
	req.Header.Add("dnt", "1")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://slickdeals.net/live/")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	res, err := m.Client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Find the review items
	data := doc.Find("#detailsDescription").Text()
	image, exists := doc.Find("#mainImage").Attr("src")
	if exists == false {
		image = "https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png'"
	}
	productlink, ex := doc.Find("#detailsTop > div > div.detailRightWrapper.forumThread > div.detailImages > div.mainImageContainer > a").Attr("href")

	if ex == false {
		//= "https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png"
		productlink, e := doc.Find("#detailsDescription > a").Attr("href")
		if e == false {
			productlink = "https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png"
		} else if !strings.Contains(productlink, "slickdeals.net/?") {
			productlink = strings.Split(productlink, "?")[0]
			
		} else {
			productlink = m.getRealLink(productlink)
		}
	} else if !strings.Contains(productlink, "slickdeals.net/?") {
		productlink = strings.Split(productlink, "?")[0]
		productlink = m.getRealLink(productlink)
	} else {

		productlink = m.getRealLink(productlink)
	}
	fmt.Println(image, productlink, exists, ex)
	fmt.Println(strings.TrimSpace(data))
	return strings.TrimSpace(data), image, productlink
}

func (m *Monitor) checkStop() error {
	for !m.stop {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
			}
		}()
		url := fmt.Sprintf("https://monitors-9ad2c-default-rtdb.firebaseio.com/monitor/%s/%s.json", strings.ToUpper(m.Config.site), m.Config.sku)
		req, _ := http.NewRequest("GET", url, nil)
		res, _ := http.DefaultClient.Do(req)

		body, _ := ioutil.ReadAll(res.Body)
		var currentObject ItemInMonitorJson
		err := json.Unmarshal(body, &currentObject)
		if err != nil {
			fmt.Println(err)

		}
		m.stop = currentObject.Stop
		m.CurrentCompanies = currentObject.Companies
		fmt.Println(m.CurrentCompanies)
		res.Body.Close()
		//	fmt.Println(currentObject)
		time.Sleep(5000 * (time.Millisecond))
	}
	return nil
}

func (m *Monitor) getRealLink(url string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Real Link Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)
	url1 := res.Request.URL.String()
	fmt.Println(url1)
	fmt.Println(strings.Split(url1, "?")[0])
	return strings.Split(url1, "?")[0]
}
