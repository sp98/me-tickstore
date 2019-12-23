package store

var (
	dailyOHLCQuery          = "SELECT * FROM %s ORDER BY time DESC"
	latestStockDetailsQuery = "SELECT last(LastPrice) as LastPrice, last(AverageTradePrice) as AverageTradePrice, last(TotalBuy) as TotalBuy, last(TotalSell) as TotalSell, last(Open) as Open, last(High) as High, last(Low) as Low, last(Close) as Close  FROM %s"
)
