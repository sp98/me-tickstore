package store

import (
	"fmt"
	"log"
	"time"

	client "github.com/orourkedd/influxdb1-client/client"
)

//DB is the influx db struct
type DB struct {
	Address     string
	Name        string
	Measurement string
}

//NewDB returns instance of an InfluxDB struct
func NewDB(address string, name string, measurement string) *DB {
	return &DB{
		Address:     address,
		Name:        name,
		Measurement: measurement,
	}

}

//GetClient creates a new Influx DB client
func (db *DB) GetClient() (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: db.Address,
	})
	if err != nil {
		log.Fatalln("Error on creating Influx DB client: ", err)
		return nil, err
	}
	return c, nil
}

func (db DB) executeQuery(query client.Query) (*client.Response, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()

	response, err := dbClient.Query(query)

	if err != nil && response.Error() != nil {
		log.Fatalln("Error executing Query - ", err)
		return nil, err
	}

	return response, nil

}

//GetDailyOHCL fetches OHLC data for today
func (db DB) GetDailyOHCL() (*client.Response, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		Command:  fmt.Sprintf(dailyOHLCQuery, db.Measurement, time.Now().Format("2006-01-02")),
		Database: db.Name,
	}
	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		return nil, err
	}
	return response, nil

}

//GetStockDetails fetches details about the stock like price, totalsetll etc.
func (db DB) GetStockDetails() (*client.Response, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		Command:  fmt.Sprintf(latestStockDetailsQuery, db.Measurement),
		Database: db.Name,
	}
	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		return nil, err
	}
	return response, nil

}
