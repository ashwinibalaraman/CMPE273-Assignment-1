package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/gorilla/sessions"
        _ "./mattn/go-yql"
)

type FinanceApiService struct{}

type ResParameters struct {
	TradeId          string
	Stocks           string
	UninvestedAmount float64
}

type ReqParameters struct {
	Symbols         string
	AmountPerSymbol [10]float64
}

var SessionResponse *ResParameters

func (h *FinanceApiService) GetQuote(r *http.Request, args *ReqParameters, Response *ResParameters) error {
	fmt.Println("---IN SERVER--")
	var data map[string]interface{}

	db, _ := sql.Open("yql", "")
	fmt.Println("Symbols:", args.Symbols)
	query := "select symbol,Ask from yahoo.finance.quotes where symbol in (" + args.Symbols + ")"
	fmt.Println("query:", query)
	stmt, err := db.Query(query)
	//db.QueryRow(query)

	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	var j int
	var numberOfStocks [10]int

	for stmt.Next() {

		stmt.Scan(&data)
		fmt.Printf("%v\n", data["symbol"])
		fmt.Printf("  %v\n\n", data["Ask"])

		var ask float64
		var symbol string

		if str1, ok := data["symbol"].(string); ok {
			symbol = str1
		}
		if str2, ok := data["Ask"].(string); ok {
			ask, _ = strconv.ParseFloat(str2, 64)
		}

		numberOfStocks[j] = int(args.AmountPerSymbol[j] / ask)
		numStockstr := strconv.Itoa(numberOfStocks[j])
		askStr := strconv.FormatFloat(ask, 'f', -1, 64)
		if Response.Stocks == "" {
			Response.Stocks = Response.Stocks + symbol + ":" + numStockstr + ":" + askStr + ","
		} else {
			Response.Stocks = Response.Stocks + symbol + ":" + numStockstr + ":" + askStr
		}

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		x := r1.Intn(1000)
		Response.TradeId = strconv.Itoa(x)
		Response.UninvestedAmount = Response.UninvestedAmount + math.Mod(args.AmountPerSymbol[j], ask)
		j++
	}

	//fmt.Println("TradeId", Response.TradeId)
	fmt.Println("Stocks:", Response.Stocks)
	fmt.Println("uninvested amount", Response.UninvestedAmount)
	//SessionResponse = Response
	//MyHandler(, r, Response)
	return nil
}

var Store = sessions.NewCookieStore([]byte("something-very-secret"))

//var Session *sessions.Session

func initSession(r *http.Request) *sessions.Session {
	Store.Options = &sessions.Options{
		MaxAge:   86400 * 7, // 1 hour
		HttpOnly: true,
	}
	session, err := Store.Get(r, "golang_cookie")
	if err != nil {
		panic(err)
	}
	return session
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	//session, err := store.Get(r, "session-name")
	session := initSession(r)

	gob.Register(&ResParameters{})
	session.Values["response"] = SessionResponse
	session.Save(r, w)
	fmt.Println("IN MY HANDLER")
	//Session = session
	//val := session.Values["response"]
	//var response ResParameters
	//response, ok := val.(*ResParameters)
	//if !ok {
	//fmt.Println("response not of type struct ResParameters")
	//}
	//fmt.Println("Stocks for tradeid", response.TradeId, " is ", response.Stocks)
}

func main() {
	RPC := rpc.NewServer()
	RPC.RegisterCodec(json.NewCodec(), "application/json")
	RPC.RegisterService(new(FinanceApiService), "")
	http.Handle("/", RPC)
	//http.HandleFunc("/rpc", MyHandler)
	//h := http.HandlerFunc(MyHandler)
	/*mux := http.NewServeMux()
	mux.Handle("/", handler(MyHandler))*/
	log.Println("Starting JSON-RPC server on localhost:1234/RPC2")
	log.Fatal(http.ListenAndServe(":1234", context.ClearHandler(http.DefaultServeMux)))
}
