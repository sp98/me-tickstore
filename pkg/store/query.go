package store

var (
	dailyOHLCQuery          = "SELECT * FROM %s  where time >= '%s' ORDER BY time DESC LIMIT 15"
	timedOHLCQuery          = "SELECT * FROM %s  where time >='%s' and time <= '%s' ORDER BY time DESC"
	latestStockDetailsQuery = "SELECT last(LastPrice) as LastPrice, last(AverageTradePrice) as AverageTradePrice, last(TotalBuyQuantity) as TotalBuy, last(TotalSellQuantity) as TotalSell, last(Open) as Open, last(High) as High, last(Low) as Low, last(Close) as Close  FROM %s"
)
