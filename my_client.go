package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/rpc/json"
)

type ResParameters struct {
	TradeId          string
	Stocks           string
	UninvestedAmount float64
}

type ReqParameters struct {
	Symbols         string
	AmountPerSymbol [10]float64
}

type CheckPortfolioResponse struct {
	stocks             string
	currentMarketValue float64
	uninvestedAmount   float64
}

//var Store = sessions.NewCookieStore([]byte("something-very-secret"))

//var Session *sessions.Session

func BuyStocks(method string, args ReqParameters) (Response ResParameters, err error) {
	buf, _ := json.EncodeClientRequest(method, args)
	fmt.Println("Symbols in jsonrpccall", args.Symbols)

	//fmt.Println("buffer : ", bytes.NewBuffer(buf))
	req, err := http.NewRequest("POST", "http://localhost:1234/rpc", bytes.NewBuffer(buf))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	/*	defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	*/
	err = json.DecodeClientResponse(resp.Body, &Response)
	/*if err != nil {
		return
	}*/
	//var store = sessions.NewCookieStore([]byte("something-very-secret"))
	/*	cookies := resp.Cookies()
		//cookie, _ := req.Cookie("golang-cookie")
		fmt.Println("cookies", cookies)
		fmt.Println("Session name", Session.Name())
	*/

	return
}

func CheckPortfolio(method string, requestParams ReqParameters) (PortfolioResponse ResParameters, err error) {

	PortfolioResponse, err = BuyStocks("FinanceApiService.GetQuote", requestParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("****portfolio response***", PortfolioResponse)

	//	ResponseMap1["TradeId"] = PortfolioResponse.TradeId
	//ResponseMap1["Stocks"] = PortfolioResponse.Stocks
	//ResponseMap1["UninvestedAmount"] = PortfolioResponse.UninvestedAmount

	return
}

var StockMap map[string][]map[string]interface{}

func manageResponseParams(Response ResParameters, TradeId string) {
	//var StockMap map[string]map[string]interface{}

	var StockInfoMap map[string]interface{}

	fmt.Println("------RESPONSE stocks", Response.Stocks)
	fmt.Println("------RESPONSE MAP", StockMap)

	stocks := strings.Split(Response.Stocks, ",")

	for _, v := range stocks {
		eachStock := strings.Split(v, ":")

		symbol := eachStock[0]
		numberStocks, _ := strconv.Atoi(eachStock[1])
		stockValue := eachStock[2]

		/*if StockMap[TradeId] != nil {
				for key, StockInfoMap := range StockMap {

					//fmt.Println("tradeid from prev session", Response.TradeId)
					//StockInfoMap = StockMap[Response.TradeId]
					fmt.Println("fetched trade id", key)
					fmt.Println("map for trade id", StockMap[key])
					if symbol == StockInfoMap["Symbol"] {
						if numStock, ok := StockInfoMap["NumStocks"].(int); ok {
							numberStocks = numStock + numberStocks
						}
					}//symbol match

				}//for loop StockMap


		}*/ // if TradId not nill
		StockInfoMap = make(map[string]interface{})
		StockInfoMap["Symbol"] = symbol
		StockInfoMap["NumStocks"] = numberStocks
		StockInfoMap["StockValue"] = stockValue
		//TempMap := StockMap["TradeId"]
		//TempMap["Symbol"] = StockInfoMap
		StockMap[TradeId] = append(StockMap[TradeId], StockInfoMap)
	} // for stocks
	//return StockMap
}

func manageCommandLineArgs(stockSymbolAndPercentage string, budget float64, ReqParams ReqParameters) (ReqParameters, error) {

	var totalPercent float64
	ReqParams.Symbols = ""
	eachStock := strings.Split(stockSymbolAndPercentage, ",")

	for i := range eachStock {
		sp := strings.Split(eachStock[i], ":")

		ReqParams.Symbols = ReqParams.Symbols + "'" + sp[0] + "'"
		if i < len(eachStock)-1 {
			ReqParams.Symbols = ReqParams.Symbols + ","
		}
		a := strings.Split(sp[1], "%")

		fmt.Println("a", a[0])
		p, _ := strconv.ParseFloat(a[0], 64)
		fmt.Println("i", i)
		ReqParams.AmountPerSymbol[i] = p / 100 * budget
		fmt.Println("amount", i, ReqParams.AmountPerSymbol[i])
		totalPercent = totalPercent + p
		i++
	}
	if totalPercent > 100 {
		var err error
		err = errors.New("Individual percents do not sum upto 100")
		fmt.Println("Individual percents do not sum upto 100")
		return ReqParams, err
	}
	fmt.Println("Symbols in func", ReqParams.Symbols)
	return ReqParams, nil
}

func main() {
	var stockSymbolAndPercentage string
	var budget float64

	StockMap = make(map[string][]map[string]interface{})
	var GetQuoteResponseMap []map[string]interface{}
	var PortfolioResponseMap []map[string]interface{}
	GetQuoteResponseMap = make([]map[string]interface{}, 10)
	PortfolioResponseMap = make([]map[string]interface{}, 10)

	ReqParams := ReqParameters{}
	ResParams := ResParameters{}
	PortfolioResponse := ResParameters{}

	var err error

	/*	flag.StringVar(&stockSymbolAndPercentage, "ssap", "GOOG:50%", "a string")
		flag.Float64Var(&budget, "budget", 2000, "a float")
		flag.Parse()
	*/

	for {
		fmt.Println("Enter 1 for Buying Stocks, 2 for Portfolio 3 to Exit\n")
		var Choice int
		fmt.Scanln(&Choice)

		if Choice == 1 {

			fmt.Println("Enter stock symbol and percentage in the format: GOOG:50%\n")
			fmt.Scanln(&stockSymbolAndPercentage)
			if budget == 0 {
				fmt.Println("Enter budget\n")
				fmt.Scanln(&budget)
			}

			ReqParams, err = manageCommandLineArgs(stockSymbolAndPercentage, budget, ReqParams)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Symbols in main", ReqParams.Symbols)
			fmt.Println("Amount in main", ReqParams.AmountPerSymbol)

			ResParams, err = BuyStocks("FinanceApiService.GetQuote", ReqParams)
			if err != nil {
				log.Fatal(err)
			}

			manageResponseParams(ResParams, ResParams.TradeId)
			GetQuoteResponseMap = StockMap[ResParams.TradeId]

			fmt.Println("\n\n\n----------BUYING STOCK SUMMARY-----------")
			fmt.Println("TradeId: ", ResParams.TradeId)
			fmt.Println("Stocks: ", ResParams.Stocks)
			fmt.Println("uninvested amount", ResParams.UninvestedAmount)
			fmt.Println("-----------------------------------------\n\n\n")

			//moneyUsed := budget - ResParams.UninvestedAmount
			//budget = budget - moneyUsed
		}

		if Choice == 2 {
			var TradeId string
			fmt.Println("Enter TradId for the portfolio you want to retrieve\n")
			fmt.Scanln(&TradeId)

			PortfolioResponse, err = CheckPortfolio("FinanceApiService.GetQuote", ReqParams)
			if err != nil {
				log.Fatal(err)
			}

			manageResponseParams(PortfolioResponse, "")
			PortfolioResponseMap = StockMap[TradeId]
			fmt.Println("GetQuoteResponseMap length: ", len(GetQuoteResponseMap))
			fmt.Println("PortfolioResponseMap length: ", len(PortfolioResponseMap))
			fmt.Println("PortfolioResponseMap ", PortfolioResponseMap)

			CPR := CheckPortfolioResponse{}

			//ArrayMapPortfolio := PortfolioResponseMap[TradeId]
			//	ArrayMapGenerate := GetQuoteResponseMap[TradeId]

			var i int
			for i < len(PortfolioResponseMap) {
				AMP := PortfolioResponseMap[i]
				AMG := GetQuoteResponseMap[i]
				var sign string
				var firstStockValue float64
				var secondStockValue float64
				if a1, ok := AMP["StockValue"].(string); ok {
					firstStockValue, _ = strconv.ParseFloat(a1, 64)
				}
				if a2, ok := AMG["StockValue"].(string); ok {
					secondStockValue, _ = strconv.ParseFloat(a2, 64)
				}
				if firstStockValue > secondStockValue {
					sign = "+"
				} else if firstStockValue < secondStockValue {
					sign = "-"
				} else {
					sign = "="
				}
				numStocks := strconv.Itoa(AMP["NumStocks"].(int))
				stockVal := AMP["StockValue"].(string)
				if i < len(PortfolioResponseMap)-1 {
					CPR.stocks = CPR.stocks + AMP["Symbol"].(string) + ":" + numStocks + sign + stockVal + ","
				} else {
					CPR.stocks = CPR.stocks + AMP["Symbol"].(string) + ":" + numStocks + ":" + sign + stockVal
				}
				CPR.currentMarketValue = CPR.currentMarketValue + float64(AMP["NumStocks"].(int))*secondStockValue
				i++
			}

			CPR.uninvestedAmount = budget - CPR.currentMarketValue
			fmt.Println("\n\n\n-------------YOUR PORTFOLIO:-------------")
			fmt.Println("STOCKS: ", CPR.stocks)
			fmt.Println("Current Market Value: ", CPR.currentMarketValue)
			fmt.Println("Uninvested Amount: ", CPR.uninvestedAmount)
			fmt.Println("-----------------------------------------\n\n\n")
			/*	session, err := Store.Get(req, "golang_cookie")
				if err != nil {
					panic(err)
				}
				//var appSession *sessions.Session
				if session == nil {
					fmt.Println("session nil")
				} else {
					fmt.Println("SESSION name", session.Name())

					val := session.Values["response"]
					//var response ResParameters
					response, ok := val.(*ResParameters)
					if !ok {
						fmt.Println("response not of type struct ResParameters")
					} else {
						fmt.Println("Stocks for tradeid", response.TradeId, " is ", response.Stocks)
					}
				}*/
		}
		if Choice == 3 {
			fmt.Println("Choice 3 entered: Shutting Down Client")
			break
		}

	}

}
