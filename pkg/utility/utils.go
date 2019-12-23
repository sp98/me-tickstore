package utility

import "strings"

//GetFilteredStocks returns formated 2-D array of stocks
func GetFilteredStocks(stock string, tokens []string) [][]string {
	result := [][]string{}
	stocks := strings.Split(stock, ",")

	for _, s := range stocks {
		sSlice := strings.Split(s, ";")
		if contains(tokens, sSlice[2]) {
			result = append(result, sSlice)
		}
	}

	return result
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
