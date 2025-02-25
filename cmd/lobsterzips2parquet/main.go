package main

import (
	"flag"

	"github.com/samber/lo"

	"github.com/hanselrd/mdviewer/internal/lobster"
)

func main() {
	output := flag.String("o", "", "output parquet file")
	flag.Parse()
	inputs := flag.Args()

	lo.Must0(lobster.ConvertZipsToParquet(inputs, *output))
}
