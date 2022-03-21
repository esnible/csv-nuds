// Produce NUDS from a numismatic CSV file

package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/esnible/csv-nuds/converter"
)

// Convert CSV to NUDS
// nolint: funlen
func main() {

	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Fprintf(os.Stderr, "syntax: %s <outputdir> <csvname> [<csvname>]\n", os.Args[0])
		os.Exit(3)
	}

	dirName := os.Args[1]
	csvName := os.Args[2]
	csvEveryName := os.Args[3]

	// We will generate one record for every row in the .CSV
	csvCoinReader, cols, err := csvReader(csvName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// If there is a second .csv, it should have only a single
	// row.  The row's values will be applied to every generated coin.
	// It is for applying stuff that should appear on every record,
	// such as the owner, copyright, database export timestamp, etc.

	var colsEveryCoin map[int]string

	var recEveryCoin []string

	if len(os.Args) == 4 {
		var csvEveryCoinReader *csv.Reader

		csvEveryCoinReader, colsEveryCoin, err = csvReader(csvEveryName)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		recEveryCoin, err = csvEveryCoinReader.Read()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	converter := converter.NewConverter(time.Now())

	// Go through each row in the CSV, producing a <NUDS> for each
	for {
		rec, err := csvCoinReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		coin := generateMap(cols, rec, colsEveryCoin, recEveryCoin)

		nuds, err := converter.GenerateNUDS(coin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fXML, err := os.Create(filepath.Join(dirName, nuds.Control.RecordID+".xml"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		encoder := xml.NewEncoder(fXML)
		encoder.Indent(" ", "  ")

		err = encoder.Encode(nuds)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

// csvReader() opens a .csv file, returning a reader and a column-to-header lookup
func csvReader(fileName string) (*csv.Reader, map[int]string, error) {
	fCSV, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	// defer fCSV.Close()

	csvReader := csv.NewReader(fCSV)

	header, err := csvReader.Read()
	if err != nil {
		return nil, nil, err
	}

	return csvReader, generateColumnLookup(header), nil
}

// generateColumnLookup() creates a map of column position to column header name
func generateColumnLookup(cols []string) map[int]string {
	retval := map[int]string{}

	for col, heading := range cols {
		retval[col] = heading
	}

	return retval
}

// generateMap creates a key=>value lookup from a row of data values
func generateMap(colLookup map[int]string, coin []string,
	everyColLookup map[int]string, every []string) map[string]string {

	retval := map[string]string{}
	applyLookup(retval, every, everyColLookup)
	applyLookup(retval, coin, colLookup)
	return retval
}

func applyLookup(coin map[string]string, vals []string, lookup map[int]string) {

	for col, val := range vals {
		// Skip if the field is unset
		if val == "" {
			continue
		}

		coin[strings.ToLower(lookup[col])] = val
	}
}
