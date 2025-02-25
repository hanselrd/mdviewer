package lobster

type Message struct {
	Time      int64     `parquet:"time,timestamp(nanosecond),delta"`
	EventType EventType `parquet:"event_type"`
	OrderID   int       `parquet:"order_id,delta"`
	Size      int       `parquet:"size"`
	Price     int       `parquet:"price,delta"`
	Side      Side      `parquet:"side"`
}

func (m Message) Status() Status {
	if m.EventType != EventTypeStatus {
		return StatusReserved
	}

	switch m.Price {
	case -1:
		return StatusHalted
	case 0:
		return StatusQuoting
	case 1:
		return StatusTrading
	}

	panic(m.Price)
}

type EventType uint

const (
	EventTypeReserved EventType = iota
	EventTypeNewLimitOrder
	EventTypeCancelLimitOrder
	EventTypeDeleteLimitOrder
	EventTypeExecuteVisibleLimitOrder
	EventTypeExecuteHiddenLimitOrder
	EventTypeCrossTrade
	EventTypeStatus
)

type Side int

const (
	SideSell Side = -1
	SideBuy       = 1
)

type Status uint

const (
	StatusReserved Status = iota
	StatusHalted
	StatusQuoting
	StatusTrading
)
