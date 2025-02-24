package lobster

type OrderBookUpdate struct {
	Symbol string `parquet:"symbol,dict"`
	Message
	OrderBook
}
