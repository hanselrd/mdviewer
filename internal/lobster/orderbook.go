package lobster

type OrderBook struct {
	Bids []PriceLevel `parquet:"bids"`
	Asks []PriceLevel `parquet:"asks"`
}

type PriceLevel struct {
	Price int `parquet:"price,delta"`
	Size  int `parquet:"size,delta"`
}

func (pl PriceLevel) Occupied() bool {
	switch pl.Price {
	case -9999999999, 9999999999:
		return pl.Size == 0
	}

	return false
}
