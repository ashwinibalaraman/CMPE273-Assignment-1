# CMPE273-Assignment-1

You will be building a virtual stock trading system for whoever wants to learn how to invest in stocks. The system must use real-time pricing via Yahoo finance API and will support USD currency only. The system has two features: Buying stocks Request “stockSymbolAndPercentage”: string (E.g. “GOOG:50%,YHOO:50%”) “budget” : float32 Response “tradeId”: number “stocks”: string (E.g. “GOOG:100:$500.25”, “YHOO:200:$31.40”) “unvestedAmount”: float32 Checking your portfolio (loss/gain) Request “tradeId”: number

Response “stocks”: string (E.g. “GOOG:100:+$520.25”, “YHOO:200:-$30.40”) “currentMarketValue” : float32 “unvestedAmount”: float32 The system will have 2 components: client and server. server: the trading engine will have JSON-RPC interface for the above features. client: the JSON-RPC client will take command line input and send requests the server.

Solution

Curl command to run curl -vX POST -H "X-Custom-Header:myvalue" -H "Content-Type:application/json" -d '{"method":"FinanceApiService.GetQuote","params":[{"Symbols":"\"GOOG\",\"YHOO\"","AmountPerSymbol":[50,50,0,0,0,0,0,0,0,0]}],"id":5577006791947779410}' http://localhost:1234/rpc

Sample curl output: Ashwini:abc rohitkandhari$ curl -vX POST -H "X-Custom-Header:myvalue" -H "Content-Type:application/json" -d '{"method":"FinanceApiService.GetQuote","params":[{"Symbols":"\"GOOG\",\"YHOO\"","AmountPerSymbol":[50,50,0,0,0,0,0,0,0,0]}],"id":5577006791947779410}' http://localhost:1234/rpc

About to connect() to localhost port 1234 (#0)
Trying ::1...
connected
Connected to localhost (::1) port 1234 (#0) > POST /rpc HTTP/1.1 > User-Agent: curl/7.27.0 > Host: localhost:1234 > Accept: / > X-Custom-Header:myvalue > Content-Type:application/json > Content-Length: 149 >
upload completely sent off: 149 out of 149 bytes < HTTP/1.1 200 OK < Content-Type: application/json; charset=utf-8 < X-Content-Type-Options: nosniff < Date: Sun, 04 Oct 2015 02:36:47 GMT < Content-Length: 125 < {"result":{"TradeId":"207","Stocks":"GOOG:0:627,YHOO:1:30.7","UninvestedAmount":69.3},"error":null,"id":5577006791947779410}
Connection #0 to host localhost left intact
Closing connection #0 Ashwini:abc rohitkandhari$
Sample client server interaction output Ashwini:abc rohitkandhari$ go run my_client.go Enter 1 for Buying Stocks, 2 for Portfolio 3 to Exit

1 Enter stock symbol and percentage in the format: GOOG:50%

GOOG:50%,YHOO:50% Enter budget

1000

----------BUYING STOCK SUMMARY----------- TradeId: 917 Stocks: GOOG:0:627,YHOO:16:30.7

uninvested amount 508.8

Enter 1 for Buying Stocks, 2 for Portfolio 3 to Exit

2 Enter TradId for the portfolio you want to retrieve

917

-------------YOUR PORTFOLIO:------------- STOCKS: GOOG:0=627,YHOO:16:=30.7 Current Market Value: 491.2

Uninvested Amount: 508.8

Enter 1 for Buying Stocks, 2 for Portfolio 3 to Exit

3 Choice 3 entered: Shutting Down Client Ashwini:abc rohitkandhari$
