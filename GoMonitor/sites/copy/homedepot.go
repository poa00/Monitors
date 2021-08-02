package Copy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/nickname32/discordhook"
	"github.com/pkg/errors"

	Webhook "github.con/prada-monitors-go/helpers/discordWebhook"
	FetchProxies "github.con/prada-monitors-go/helpers/proxy"
	Types "github.con/prada-monitors-go/helpers/types"
)

type CurrentMonitor struct {
	Monitor Types.Monitor
}

type Payload struct {
	OperationName string    `json:"operationName"`
	Variables     Variables `json:"variables"`
	Query         string    `json:"query"`
}

type Variables struct {
	ExcludeInventory bool     `json:"excludeInventory"`
	ItemIds          []string `json:"itemIds"`
	StoreID          string   `json:"storeId"`
}


func NewMonitor(sku string) *CurrentMonitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("TESTING")
	m := CurrentMonitor{}
	m.Monitor.Availability = "Copy"
	m.Monitor.Config.Site = "Copy"
	m.Monitor.Config.Sku = sku
	m.Monitor.Client = http.Client{Timeout: 10 * time.Second}
	proxyList := FetchProxies.Get()

	// time.Sleep(15000 * (time.Millisecond))
	// go m.Monitor.CheckStop()
	// time.Sleep(3000 * (time.Millisecond))

	i := true
	for i {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Monitor.Config.Site, m.Monitor.Config.Sku, r)
			}
		}()
		if !m.Monitor.Stop {
			currentProxy := m.Monitor.GetProxy(proxyList)
			splittedProxy := strings.Split(currentProxy, ":")
			proxy := Types.Proxy{splittedProxy[0], splittedProxy[1], splittedProxy[2], splittedProxy[3]}
			//	fmt.Println(proxy, proxy.ip)
			prox1y := fmt.Sprintf("http://%s:%s@%s:%s", proxy.UserAuth, proxy.UserPass, proxy.Ip, proxy.Port)
			proxyUrl, err := url.Parse(prox1y)
			if err != nil {
				fmt.Println(errors.Cause(err))
				return nil
			}
			defaultTransport := &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
			m.Monitor.Client.Transport = defaultTransport
			go m.monitor()
			time.Sleep(500 * (time.Millisecond))
		} else {
			fmt.Println(m.Monitor.Config.Sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *CurrentMonitor) monitor() error {
	// watch := stopwatch.Start()
	// // Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	// // curl 'https://www.homedepot.com/product-information/model?opname=mediaPriceInventory' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' --compressed -H 'Referer: https://www.homedepot.com/p/Milwaukee-M18-18-Volt-Lithium-Ion-Cordless-Combo-Tool-Kit-7-Tool-with-Two-3-0-Ah-Batteries-Charger-and-Tool-Bag-2695-27S/312422229' -H 'content-type: application/json' -H 'X-Experience-Name: general-merchandise' -H 'apollographql-client-name: general-merchandise' -H 'apollographql-client-version: 0.0.0' -H 'X-current-url: /p/Milwaukee-M18-18-Volt-Lithium-Ion-Cordless-Combo-Tool-Kit-7-Tool-with-Two-3-0-Ah-Batteries-Charger-and-Tool-Bag-2695-27S/312422229' -H 'X-Cloud-Trace-Context: 0e64ce7a7684f25a0a0b14a36629a3bd/5838804585739945946' -H 'X-Api-Cookies: {"x-user-id":"7d5a250a-8cfa-d513-e384-ee321345ad5c"}' -H 'x-debug: false' -H 'Origin: https://www.homedepot.com' -H 'Connection: keep-alive' -H 'Cookie: RT="z=1&dm=www.homedepot.com&si=4944bb16-147e-45e4-91c9-9039983ca8f6&ss=kr8pdzz8&sl=8&tt=b1d&bcn=%2F%2F3211c0e1.akstat.io%2F&ld=6hvt"; _pxde=44bd837dbe40aeb196a7b605be8bcede4f37e7cf667f6d5e512b88987a7342b1:eyJ0aW1lc3RhbXAiOjE2MjY1ODMwODAxMDZ9; AKA_A2=A; HD_DC=origin; _abck=8506091DCBFB383148981DF6EE2A7D5A~0~YAAQvRY2F1wqbKN6AQAAjiHotwYdx3jY1lIFa+yfqOwsAGGw7OXfixfnVCo/jhxxkb9I2xrp0PuYKoJRBbBSYSEDFKaD0gaA2Bm6rS6jknuoI71pE+kKkWVNp+rVV4XVavnCcBvLcvrDHMLW9/43Gajr2M8gB5WMgGU0Sr49o04G5MHLjQqLGilPHsrQuk33woxyINwhbS9KWeqgj+7r1ZLV2VtnpzWJSN0+Bw7v3NuQRPq7OVS/2TIcRFdtq+jQUv5Io+6KyfRCOqPX6dQxg+jdx8xZXgtY17fZfq/U2I4VVvcENbKPn+NC8pi+lO40Eop/Ch15MJaXrdcR15c0PiS4tJQUz78f+nq/NSuF8GaSRQ03pPf7O6ne9nmHd8BsLV3yYjedu2/dm23hedBRHLagZEpBOZpBd+PYCg==~-1~-1~-1; ak_bmsc=EE82CE4BAD1DA787F660C14A2E9A8EC3~000000000000000000000000000000~YAAQlRY2F8wALK56AQAAPnbltwxO+fOBfgmDZ6OXrC4QQLB5ySICgN5RaK7vC9+SaxMtI862a4mMUg+D6+DmuQu2XXut8j/BhBI2zijOvSKFhtk2wJ43rmFizS+44oO5Eeh/WmCO+kQ9RUtl3riPndqr3eVzu37xJ/MZWsu+50n/zQm2Aa7t2X0mD5evx4h+0znECJvk//AXEMNt9ShUWT7DNfuUvuaKSFjZyEE4WeEzA5ZGd4lgk+3blkeIsu3whEbUZ2dmwkqmpdeuqWJyiacTiKQ/Z9pzuxpABuN4KsLe2wwHh9gQJR0lFaweqhL1UjgXJziypXaSOoNWKfeZ5kRXp3ApnAJErU0U8IFXu4/gU9W89PHFv+Ts4u2r1IvgKKYZot/QZJZx/H9EqkzWLH/yfMeqXJrBmj83c19nwiHuWH8f4I0u5NEybTUw9Zvkk8xFni0m/3G3Fzr/0vfgn6eM929FFblJsnP3g6gb1o3ccATJSqvlmL1lTl5a; bm_sz=855D66AF60C8B0560C2A97081A5889C0~YAAQlRY2F8YALK56AQAAUnHltww5guNs519d7b0zVeyadNLiToLSS5cJLLRoVwt50qrQLp7+Hi9yAfLyVS4Sa97z5odOeNAvh8K0X1c1bIzxHHwV4ESLNeGIAQKo6Dm6Q1u/c0lzhA19QONvmb18fyi+U/K7BlGk6xX7UBLk0HoOiqzShhyHNRkmgYnkqCGKhOEt8io8GPRBmIxWkqIY7KuWNVBXSsmjQneJLGTmpyr2of6eBkpYiqhenVPy9g4lVk2fkGrlbWUfTPsBv2tSufz7A4OH4JnQ98bgk87RxGUTdRO7wDo=~3422532~3748678; THD_NR=1; AMCV_F6421253512D2C100A490D45%40AdobeOrg=1585540135%7CMCIDTS%7C18827%7CMCMID%7C53283005719650550872622373288851036139%7CMCAAMLH-1627187707%7C9%7CMCAAMB-1627187707%7C6G1ynYcLPuiQxYZrsz_pkqfLG9yMXBpb2zX5dvJdYQJzPXImdj0y%7CMCOPTOUT-1626590107s%7CNONE%7CMCCIDH%7C-1161900227%7CvVersion%7C4.4.0; at_check=true; mbox=session#0d32cbcc779f44bba05d970d2eda9a09#1626584766|PC#0d32cbcc779f44bba05d970d2eda9a09.34_0#1689827881; THD_PERSIST=C4%3d224%2bBoynton%20Beach%20%2d%20Boynton%20Beach%2c%20FL%2b%3a%3bC4%5fEXP%3d1678422973%3a%3bC6%3d%7b%22I1%22%3a%221%22%7d%3a%3bC6%5fEXP%3d1629174973%3a%3bC24%3d33426%3a%3bC24%5fEXP%3d1678422973%3a%3bC39%3d1%3b7%3a00%2d20%3a00%3b2%3b6%3a00%2d22%3a00%3b3%3b6%3a00%2d22%3a00%3b4%3b6%3a00%2d22%3a00%3b5%3b6%3a00%2d22%3a00%3b6%3b6%3a00%2d22%3a00%3b7%3b6%3a00%2d22%3a00%3a%3bC39%5fEXP%3d1648182973; THD_SESSION=; THD_CACHE_NAV_SESSION=; THD_CACHE_NAV_PERSIST=; IN_STORE_API_SESSION=TRUE; WORKFLOW=GEO_LOCATION; THD_FORCE_LOC=1; THD_INTERNAL=0; THD_LOCALIZER=%7B%22WORKFLOW%22%3A%22GEO_LOCATION%22%2C%22THD_FORCE_LOC%22%3A%221%22%2C%22THD_INTERNAL%22%3A%220%22%2C%22THD_STRFINDERZIP%22%3A%2233426%22%2C%22THD_LOCSTORE%22%3A%22224%2BBoynton%20Beach%20-%20Boynton%20Beach%2C%20FL%2B%22%2C%22THD_STORE_HOURS%22%3A%221%3B7%3A00-20%3A00%3B2%3B6%3A00-22%3A00%3B3%3B6%3A00-22%3A00%3B4%3B6%3A00-22%3A00%3B5%3B6%3A00-22%3A00%3B6%3B6%3A00-22%3A00%3B7%3B6%3A00-22%3A00%22%2C%22THD_STORE_HOURS_EXPIRY%22%3A1626586506%7D; DELIVERY_ZIP=33426; DELIVERY_ZIP_TYPE=DEFAULT; _px_f394gi7Fvmc43dfg_user_id=ODZkNjIyNzAtZTc4MS0xMWViLWE2NWItZGRmMTg2ZjllMDE4; _pxvid=867073c7-e781-11eb-9352-0242ac120007; AMCVS_F6421253512D2C100A490D45%40AdobeOrg=1; thda.s=36a6bf4d-a476-2242-cba4-381849665f76; thda.u=7d5a250a-8cfa-d513-e384-ee321345ad5c; _px=QzjPuuYrfJ0FlMVErGQKV3TgZ+bcWiOO8Em0Tj8xVw18VhoifCdhA4wVNWKZCmZpI0SUIXJizrWIxGyICY5XGA==:1000:czkHC5XYmnX2kTQbdnYBwsvcmOUNDuzobtnz6lLKUcqes5owPrYTiLG+MUOAX3bxugE6H++ePWG6DlxczELmywbR181gQVudeu0BlHfG3Ndr5TSxVtskh4/jcXBgu9sCjOPXYYuCE8G9z2LXaPX2F6s+7cC6qsjnjB2Ipc9W+whH+N2GrsqtdbHQ+9gT9lGojGHC/SX1rRCE11Peba6u2yss5G5+zVIEhMGZrpkPutb/Atv3cMIM22dKFYThbwLOaWI1+obXbWIYPe5gxyRk2Q==; s_sess=%20s_cc%3Dtrue%3B%20s_pv_pName%3Dcart%2520add%3B%20s_pv_pType%3Dcart%2520add%2520overlay%3B%20s_pv_cmpgn%3D%3B%20s_pv_pVer%3Dmycart%2520drawer%2520v1%3B%20stsh%3D%3B; s_pers=%20s_nr365%3D1626582973891-New%7C1658118973891%3B%20s_dslv%3D1626582973891%7C1721190973891%3B; thda.m=53283005719650550872622373288851036139; forterToken=8bf956099c3846e4a3f6b14e865c4b9f_1626582924974__UDF43_13ck; ajs_user_id=null; ajs_group_id=null; ajs_anonymous_id=%228b2dec69-5c9a-449f-9a50-275af5af4b58%22; _meta_mediaMath_iframe_counter=3; _meta_bing_beaconFired=true; _meta_facebookPixel_beaconFired=true; _meta_movableInk_mi_u=8b2dec69-5c9a-449f-9a50-275af5af4b58; _meta_metarouter_timezone_offset=240; _meta_neustar_mcvisid=53283005719650550872622373288851036139; _meta_neustar_aam=53026210853881974022595153617472453022; _meta_googleFloodlight_beaconFired=true; _gcl_au=1.1.294633985.1626582908; aam_mobile=seg%3D1131078; aam_uuid=53026210853881974022595153617472453022; _meta_pinterest_epik_failure=true; _meta_revjet_revjet_vid=4960741387053710000; _meta_amobee_uid=3446291930236124480; _meta_acuityAds_auid=592924889202; _meta_acuityAds_cauid=auid=592924889202; QSI_HistorySession=https%3A%2F%2Fwww.homedepot.com%2Fp%2FHome-Decorators-Collection-Turkish-Cotton-Ultra-Soft-18-Piece-Towel-Set-in-White-NHV-8-0618WH18V%2F309947410~1626582908101%7Chttps%3A%2F%2Fwww.homedepot.com%2F~1626582914072%7Chttps%3A%2F%2Fwww.homedepot.com%2Fb%2FTools%2FSpecial-Values%2FFree-Shipping%2FN-5yc1vZc1xyZ7Z1z0zdh8%3FNCNI-5~1626582925342; _ga=GA1.2.772932008.1626582909; _gid=GA1.2.1580706087.1626582909; _meta_mediaMath_mm_id=e37960f3-af7c-4500-b017-a265eb6d0055; _meta_mediaMath_cid=e37960f3-af7c-4500-b017-a265eb6d0055; _meta_xandr_uid=0; _meta_xandr_uid2=uuid2=0; QuantumMetricSessionID=5ef8e15143b2d128e4bdd2e1a9f60610; QuantumMetricUserID=868b7ddd026e20262448e1936f90d16e; LPVID=Y3ZDRhOTllYzM5NjQ2Mjk1; LPSID-31564604=F0Xvtj2wTbONIBHblef_dQ; s_sq=homedepotglobalproduction%3D%2526c.%2526a.%2526activitymap.%2526page%253Dcart%252520add%2526link%253DShip%252520to%252520Home%252520Get%252520it%252520by%252520Thu%25252C%252520Jul%25252022%252520FREE%2526region%253Droot%2526pageIDType%253D1%2526.activitymap%2526.a%2526.c%2526pid%253Dcart%252520add%2526pidt%253D1%2526oid%253Dfunctionuc%252528%252529%25257B%25257D%2526oidt%253D2%2526ot%253DA; akaau=1626583379~id=cf4083d6b5f88c4f87d8ff194fd3adeb; _meta_adobe_aam_uuid=53026210853881974022595153617472453022; _meta_google_ga=772932008.1626582909; _meta_google_cookie_ga=GA1.2.772932008.1626582909; _meta_google_cookie_gid=GA1.2.1580706087.1626582909; trx=4993242953219991324; ads=e9234d6c7c0aaad1e9484ed1f700cc8b; THD_MCC_ID=e9d684e6-4e00-41a0-88f7-fb7d6981235d; cart_activity=b91c5337-ae56-4ba3-993a-406f2c36a1aa; RES_TRACKINGID=12871243547928562; ResonanceSegment=; RES_SESSIONID=85687343547928562' -H 'Pragma: no-cache' -H 'Cache-Control: no-cache' -H 'TE: Trailers' --data-raw '{"operationName":"mediaPriceInventory","variables":{"excludeInventory":false,"itemIds":["312422229"],"storeId":"224"},"query":"query mediaPriceInventory($excludeInventory: Boolean = false, $itemIds: [String!]!, $storeId: String!) {\n  mediaPriceInventory(itemIds: $itemIds, storeId: $storeId) {\n    productDetailsList {\n      itemId\n      imageLocation\n      onlineInventory @skip(if: $excludeInventory) {\n        enableItem\n        totalQuantity\n        __typename\n      }\n      pricing {\n        value\n        original\n        message\n        mapAboveOriginalPrice\n        __typename\n      }\n      storeInventory @skip(if: $excludeInventory) {\n        enableItem\n        totalQuantity\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n"}'
	
	// var sku []string
	// sku = append(sku, m.Monitor.Config.Sku)
	// data := Payload{
	// 	OperationName: "mediaPriceInventory",
	// 	Query: "query mediaPriceInventory($excludeInventory: Boolean = false, $itemIds: [String!]!, $storeId: String!) {  mediaPriceInventory(itemIds: $itemIds, storeId: $storeId) {    productDetailsList {      itemId      imageLocation      onlineInventory @skip(if: $excludeInventory) {        enableItem        totalQuantity        __typename      }      pricing {        value        original        message        mapAboveOriginalPrice        __typename      }      storeInventory @skip(if: $excludeInventory) {        enableItem        totalQuantity        __typename      }      __typename    }    __typename  }}",
	// 	Variables: Variables{ExcludeInventory: false, ItemIds: sku, StoreID: "224"},
	// }
	// payloadBytes, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println(errors.Cause(err))
	// 	return nil
	// }
	// payload := bytes.NewReader(payloadBytes)

	// req, err := http.NewRequest("POST", "https://www.homedepot.com/product-information/model?opname=mediaPriceInventory", payload)
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println(errors.Cause(err))
	// 	return nil
	// }
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0")
	// req.Header.Add("Accept", "*/*")
	// req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Add("Referer", "https://www.homedepot.com/p/Milwaukee-M18-18-Volt-Lithium-Ion-Cordless-Combo-Tool-Kit-7-Tool-with-Two-3-0-Ah-Batteries-Charger-and-Tool-Bag-2695-27S/312422229")
	// req.Header.Add("content-type", "application/json")
	// req.Header.Add("X-Experience-Name", "general-merchandise")
	// req.Header.Add("apollographql-client-name", "general-merchandise")
	// req.Header.Add("apollographql-client-version", "0.0.0")
	// req.Header.Add("x-debug", "false")
	// req.Header.Add("Origin", "https://www.homedepot.com")
	// req.Header.Add("Connection", "keep-alive")
	// req.Header.Add("Pragma", "no-cache")
	// req.Header.Add("Cache-Control", "no-cache")
	// req.Header.Add("TE", "Trailers")

	// res, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println(errors.Cause(err))
	// 	return nil
	// }
	// defer res.Body.Close()
	// defer func() {
	// 	watch.Stop()
	// 	fmt.Printf("Home Depot %s - Code : %d Milli elapsed: %v\n", m.Monitor.Config.Sku, res.StatusCode, watch.Milliseconds())
	// }()
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println(errors.Cause(err))
	// 	return nil
	// }
	// if res.StatusCode != 200 {
	// 	time.Sleep(10 * time.Second)
	// 	return nil
	// }

	// var jsonResponse HomeDepotResponse

	// err = json.Unmarshal([]byte(body), &jsonResponse)
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println(errors.Cause(err))
	// 	return nil
	// }
	// productQuantity := jsonResponse.Data.MediaPriceInventory.ProductDetailsList[0].OnlineInventory.TotalQuantity
	// Price := jsonResponse.Data.MediaPriceInventory.ProductDetailsList[0].Pricing.Value
	// Image := jsonResponse.Data.MediaPriceInventory.ProductDetailsList[0].ImageLocation
	// Link := fmt.Sprintf("https://www.homedepot.com/p/prada/%s", m.Monitor.Config.Sku)
	// // productName := strings.Split(jsonResponse.Data.MediaPriceInventory.ProductDetailsList[0].ImageLocation, "")
	// if productQuantity > 0 {
	// 	m.Monitor.CurrentAvailabilityBool = true
	// } else {
	// 	m.Monitor.CurrentAvailabilityBool = false
	// }


	// if !m.Monitor.AvailabilityBool && m.Monitor.CurrentAvailabilityBool {
	// 	fmt.Println("In Stock: ", productQuantity)
	// 	go m.sendWebhook(m.Monitor.Config.Sku, "Test", Price, Link, Image, productQuantity)
	// }


	// m.Monitor.AvailabilityBool = m.Monitor.CurrentAvailabilityBool
	return nil
}

func (m *CurrentMonitor) sendWebhook(sku string, name string, price float64, link string, image string, Qty int) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Monitor.Config.Site, m.Monitor.Config.Sku, r)
		}
	}()
	for _, letter := range m.Monitor.MonitorProduct.Name {
		if string(letter) == `"` {
			m.Monitor.MonitorProduct.Name = strings.Replace(m.Monitor.MonitorProduct.Name, `"`, "", -1)
		}
	}
	fmt.Println("Testing Here : ", m.Monitor.MonitorProduct.Name, "Here")
	if strings.HasSuffix(m.Monitor.MonitorProduct.Name, "                       ") {
		m.Monitor.MonitorProduct.Name = strings.Replace(m.Monitor.MonitorProduct.Name, "                       ", "", -1)
	}
	fmt.Println("Testing Here : ", m.Monitor.MonitorProduct.Name, "Here")
	// now := time.Now()
	// currentTime := strings.Split(now.String(), "-0400")[0]
	t := time.Now().UTC()

	for _, comp := range m.Monitor.CurrentCompanies {
		fmt.Println(comp.Company)
		go m.webHookSend(comp, m.Monitor.Config.Site, name, price, link, t, image, Qty)
	}
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	return nil
}

func (m *CurrentMonitor) webHookSend(c Types.Company, site string, name string, price float64, link string, currentTime time.Time, image string, Qty int) {
	Color, err := strconv.Atoi(c.Color)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return
	}
	Quantity := strconv.Itoa(Qty)
	Price := fmt.Sprintf("%f", price)
	var currentFields []*discordhook.EmbedField
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:   "Product Name",
		Value:  name,
		Inline: false,
	})
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:   "Price",
		Value:  Price,
		Inline: true,
	})
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:   "Stock #",
		Value:  Quantity,
		Inline: true,
	})
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:  "Links",
		Value: "[In Development](" + link + ")",
	})
	var discordParams discordhook.WebhookExecuteParams = discordhook.WebhookExecuteParams{
		Content: "",
		Embeds: []*discordhook.Embed{
			{
				Title:  site + " Monitor",
				URL:    link,
				Color:  Color,
				Fields: currentFields,
				Footer: &discordhook.EmbedFooter{
					Text: "Prada#4873",
				},
				Timestamp: &currentTime,
				Thumbnail: &discordhook.EmbedThumbnail{
					URL: image,
				},
				Provider: &discordhook.EmbedProvider{
					URL: c.CompanyImage,
				},
			},
		},
	}
	go Webhook.SendWebhook(c.Webhook, &discordParams)
}

