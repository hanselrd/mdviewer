package main

import (
	"github.com/samber/lo"

	"github.com/hanselrd/mdviewer/internal/lobster"
)

func main() {
	lo.Must0(
		lobster.ConvertZipsToParquet(
			[]string{
				"data/LOBSTER_SampleFile_AAPL_2012-06-21_1.zip",
				"data/LOBSTER_SampleFile_AAPL_2012-06-21_5.zip",
				"data/LOBSTER_SampleFile_AAPL_2012-06-21_10.zip",
				"data/LOBSTER_SampleFile_AAPL_2012-06-21_30.zip",
				"data/LOBSTER_SampleFile_AAPL_2012-06-21_50.zip",
				"data/LOBSTER_SampleFile_AMZN_2012-06-21_1.zip",
				"data/LOBSTER_SampleFile_AMZN_2012-06-21_5.zip",
				"data/LOBSTER_SampleFile_AMZN_2012-06-21_10.zip",
				"data/LOBSTER_SampleFile_GOOG_2012-06-21_1.zip",
				"data/LOBSTER_SampleFile_GOOG_2012-06-21_5.zip",
				"data/LOBSTER_SampleFile_GOOG_2012-06-21_10.zip",
				"data/LOBSTER_SampleFile_INTC_2012-06-21_1.zip",
				"data/LOBSTER_SampleFile_INTC_2012-06-21_5.zip",
				"data/LOBSTER_SampleFile_INTC_2012-06-21_10.zip",
				"data/LOBSTER_SampleFile_MSFT_2012-06-21_1.zip",
				"data/LOBSTER_SampleFile_MSFT_2012-06-21_5.zip",
				"data/LOBSTER_SampleFile_MSFT_2012-06-21_10.zip",
				"data/LOBSTER_SampleFile_MSFT_2012-06-21_30.zip",
				"data/LOBSTER_SampleFile_MSFT_2012-06-21_50.zip",
				"data/LOBSTER_SampleFile_SPY_2012-06-21_30.zip",
				"data/LOBSTER_SampleFile_SPY_2012-06-21_50.zip",
			},
			"data/LOBSTER.parquet",
		),
	)
}
