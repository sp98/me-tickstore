package ohlc

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sp98/tickstore/pkg/store"
)

var (
	//DBUrl is the connection url for influx db
	DBUrl = ""
	//DBName is the database name
	DBName = ""
)

//OHLC is the open, high, low and close prices
type OHLC struct {
	Open  float64 `json:"Open"`
	High  float64 `json:"High"`
	Low   float64 `json:"Low"`
	Close float64 `json:"Close"`
}

func init() {
	DBUrl = os.Getenv("INFLUX_DB_URL")
	DBName = os.Getenv("TICK_STORE_DB_NAME")
}

//Routes define the OHCL routes
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{token}/{interval}", GetOHLC)
	return router
}

//GetOHLC returns open, high, low and close price for an intrument for a particular interval
func GetOHLC(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	interval := chi.URLParam(r, "interval")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	ohlc := getTicks(token, interval, from, to)

	render.JSON(w, r, ohlc) // A chi router helper for serializing and returning json
}

func getTicks(token, interval, from, to string) []OHLC {

	db := store.NewDB(DBUrl, DBName, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", token, interval)
	response, _ := db.GetDailyOHCL(from, to)
	var ohlcList []OHLC
	for _, results := range response.Results {
		for _, rows := range results.Series {
			ohlc := &OHLC{}
			for _, row := range rows.Values {
				ohlc.Close, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[1]), 64)
				ohlc.High, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[2]), 64)
				ohlc.Low, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[3]), 64)
				ohlc.Open, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[4]), 64)
				ohlcList = append(ohlcList, *ohlc)
			}
		}

	}
	return ohlcList

}
