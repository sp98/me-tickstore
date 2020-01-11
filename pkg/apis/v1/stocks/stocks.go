package stocks

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sp98/tickstore/pkg/store"
	"github.com/sp98/tickstore/pkg/utility"
)

var (
	//DBUrl is the connection url for influx db
	DBUrl = ""
	//DBName is the database name
	DBName = ""
)

func init() {
	DBUrl = os.Getenv("INFLUX_DB_URL")
	DBName = os.Getenv("TICK_STORE_DB_NAME")
}

const (
	//StocksEnv is the enviornment variable to get all the stocks.
	StocksEnv = "STOCKS"
)

//Result stores the final response for all the stocks
type Result map[string]StockDetails

//StockDetails has fields that show stock related data
type StockDetails struct {
	Name              string  `json:"Name"`
	Symbol            string  `json:"Symnbol"`
	Token             string  `json:"Token"`
	Exchange          string  `json:"Exchange"`
	LastPrice         float64 `json:"LastPrice"`
	AverageTradePrice float64 `json:"AverageTradePrice"`
	TotalBuy          float64 `json:"TotalBuy"`
	TotalSell         float64 `json:"TotalSell"`
	Open              float64 `json:"Open"`
	High              float64 `json:"High"`
	Low               float64 `json:"Low"`
	Close             float64 `json:"Close"`
}

//Routes define the OHCL routes
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{tokens}", GetStockDetails)
	return router
}

//GetStockDetails fetch relevant stock details.
func GetStockDetails(w http.ResponseWriter, r *http.Request) {
	//Get all the stocks.
	tokens := chi.URLParam(r, "tokens")
	stocks := utility.GetFilteredStocks(os.Getenv(StocksEnv), strings.Split(tokens, "-"))
	result := Result{}
	//result := Result{}
	for _, stock := range stocks {
		token := stock[2]
		db := store.NewDB(DBUrl, DBName, "")
		db.Measurement = fmt.Sprintf("%s_%s", "ticks", token)
		// Order of Response - LastPrice TotalBuy TotalSell Open High Low Close
		response, _ := db.GetStockDetails()
		sd := StockDetails{}
		for _, results := range response.Results {
			for _, rows := range results.Series {

				for _, row := range rows.Values {
					sd.LastPrice, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[1]), 64)
					sd.AverageTradePrice, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[2]), 64)
					sd.TotalBuy, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[3]), 64)
					sd.TotalSell, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[4]), 64)
					sd.Open, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[5]), 64)
					sd.High, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[6]), 64)
					sd.Low, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[7]), 64)
					sd.Close, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[8]), 64)
					sd.Name = stock[0]
					sd.Symbol = stock[1]
					sd.Token = stock[2]
					sd.Exchange = stock[3]
				}
			}

		}

		result[token] = sd
	}

	w.Header().Set("Access-Control-Allow-Origin", "https://marketmoz.com")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, Accept-Language, Content-Length, Accept-Encoding, X-CSRF-Token")

	render.JSON(w, r, result) // A chi router helper for serializing and returning json
}
