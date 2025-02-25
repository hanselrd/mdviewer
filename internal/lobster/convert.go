package lobster

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/parquet-go/parquet-go"
	"github.com/parquet-go/parquet-go/compress/snappy"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/hanselrd/mdviewer/internal/build"
)

func ConvertZipsToParquet(ins []string, out string) error {
	for _, in := range ins {
		if filepath.Ext(in) != ".zip" {
			return fmt.Errorf("wrong extension: %s, expected: .zip", in)
		}
	}

	if filepath.Ext(out) != ".parquet" {
		return fmt.Errorf("wrong extension: %s, expected: .parquet", out)
	}

	chs := make([]chan OrderBookUpdate, len(ins))
	var eg errgroup.Group

	for i, in := range ins {
		chs[i] = make(chan OrderBookUpdate)

		eg.Go(func() error {
			return readZip(in, chs[i])
		})
	}

	eg.Go(func() error {
		return writeParquet(out, chs)
	})

	return eg.Wait()
}

func writeParquet(name string, chs []chan OrderBookUpdate) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := parquet.NewGenericWriter[OrderBookUpdate](
		file,
		parquet.CreatedBy("mdviewer", build.Version, build.Hash),
		parquet.Compression(&snappy.Codec{}),
		parquet.ColumnPageBuffers(
			parquet.NewFileBufferPool(os.TempDir(), "mdviewer.buffers.*"),
		),
	)
	defer writer.Close()

	obus := make([]OrderBookUpdate, len(chs))

	for i, ch := range chs {
		obus[i] = <-ch
	}

	for {
		obu, i := lo.MinIndexBy(obus, func(a, b OrderBookUpdate) bool {
			return a.Time < b.Time
		})

		if obu.Time == math.MaxInt64 {
			break
		}

		lo.Must(writer.Write([]OrderBookUpdate{obu}))

		if obu, ok := <-chs[i]; ok {
			obus[i] = obu
		} else {
			obus[i].Time = math.MaxInt64
		}
	}

	return nil
}

func readZip(name string, ch chan<- OrderBookUpdate) error {
	reader, err := zip.OpenReader(name)
	if err != nil {
		close(ch)
		return err
	}
	defer reader.Close()

	msgCh := make(chan []string)
	obCh := make(chan []string)
	var eg errgroup.Group

	for _, file := range reader.File {
		switch {
		case strings.Contains(file.Name, "message"):
			eg.Go(func() error {
				return readCsv(file, msgCh)
			})
		case strings.Contains(file.Name, "orderbook"):
			eg.Go(func() error {
				return readCsv(file, obCh)
			})
		}
	}

	words := lo.Words(filepath.Base(name))
	midnight := time.Date(
		lo.Must(strconv.Atoi(words[4])),
		time.Month(lo.Must(strconv.Atoi(words[5]))),
		lo.Must(strconv.Atoi(words[6])),
		0,
		0,
		0,
		0,
		lo.Must(time.LoadLocation("America/New_York")),
	)

	for msg := range msgCh {
		ob := <-obCh
		asks, bids := lo.FilterReject(
			lo.Chunk(ob, 2),
			func(_ []string, i int) bool {
				return i%2 == 0
			},
		)

		ts := strings.Split(msg[0], ".")
		if len(ts[1]) > 9 {
			ts[1] = ts[1][:9]
		}

		ch <- OrderBookUpdate{
			Symbol: fmt.Sprintf("%s:%s", words[3], words[7]),
			Message: Message{
				Time: midnight.
					Add(time.Second *
						time.Duration(lo.Must(strconv.Atoi(ts[0])))).
					Add(time.Nanosecond * time.Duration(lo.Must(strconv.Atoi(ts[1]))*int(math.Pow10(9-len(ts[1]))))).
					UnixNano(),
				EventType: EventType(lo.Must(strconv.Atoi(msg[1]))),
				OrderID:   lo.Must(strconv.Atoi(msg[2])),
				Size:      lo.Must(strconv.Atoi(msg[3])),
				Price:     lo.Must(strconv.Atoi(msg[4])),
				Side:      Side(lo.Must(strconv.Atoi(msg[5]))),
			},
			OrderBook: OrderBook{
				Bids: lo.Map(bids, func(pl []string, _ int) PriceLevel {
					return PriceLevel{
						Price: lo.Must(strconv.Atoi(pl[0])),
						Size:  lo.Must(strconv.Atoi(pl[1])),
					}
				}),
				Asks: lo.Map(asks, func(pl []string, _ int) PriceLevel {
					return PriceLevel{
						Price: lo.Must(strconv.Atoi(pl[0])),
						Size:  lo.Must(strconv.Atoi(pl[1])),
					}
				}),
			},
		}
	}
	close(ch)

	return eg.Wait()
}

func readCsv(file *zip.File, ch chan<- []string) error {
	file2, err := file.Open()
	if err != nil {
		close(ch)
		return err
	}
	defer file2.Close()

	reader := csv.NewReader(file2)

	for {
		record, err := reader.Read()
		if err != nil {
			close(ch)

			if err == io.EOF {
				return nil
			}

			return err
		}
		ch <- record
	}
}
