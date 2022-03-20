// Produce NUDS from a numismatic CSV file

package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/esnible/csv-nuds/simplenuds"
)

const (
	// Column names in CSV file we hope to support in v0.1.
	// These *MUST* be lower-case here.  In the .CSV they can be any case.
	URLCoin           = "url"
	CoinID            = "id"
	URLCoinImage      = "imageurl"
	Reporter          = "reporter"
	URLReporter       = "reporterurl"
	Category          = "category"
	Denomination      = "denomination"
	Keywords          = "keywords"
	Metal             = "metal"
	Diameter          = "diameter"
	Title             = "title"
	Weight            = "weight"
	Mint              = "mint"
	CreationTime      = "creationtime"
	AdditionalDetails = "additionaldetails"
	URLRights         = "rightsurl"
	Source            = "source"
	Date              = "date"
)

type handler func(coin *simplenuds.NUDS, val string) error

var (
	// Handlers for the different column names
	handlers = map[string]handler{
		CoinID:            recordID,
		URLCoinImage:      coinSingleURLImageHandler,
		Denomination:      denominationHandler,
		Metal:             metalHandler,
		Diameter:          diameterInMMHandler,
		Title:             titleHandler,
		Weight:            weightHandler,
		Mint:              mintHandler,
		URLRights:         rightsURLHandler,
		Source:            sourceHandler,
		CreationTime:      recordCreatedDateHandler,
		Reporter:          reporterHandler,
		AdditionalDetails: detailsHandler,
		// TODO DateHandler.  The particular dataset I used for testing
		// had 100% invalid data for date: "?", "BBA" (a mint!), and "x2"
	}
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

		nuds, err := generateNUDS(cols, rec, colsEveryCoin, recEveryCoin)
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

// generateNUDS() generates NUDS from a slice of column values (a CSV coin row) and optional second row
func generateNUDS(colLookup map[int]string, coin []string,
	everyColLookup map[int]string, every []string) (simplenuds.NUDS, error) {
	retval := simplenuds.NewNUDS("physical")

	err := applyNuds(&retval, everyColLookup, every)
	if err != nil {
		return retval, err
	}

	err = applyNuds(&retval, colLookup, coin)
	if err != nil {
		return retval, err
	}

	return retval, nil
}

// applyNuds() modifies an existing NUDS by using handlers to convert CSV values to NUDS format
func applyNuds(coin *simplenuds.NUDS, colLookup map[int]string, colVals []string) error {
	for row, val := range colVals {
		// Skip if the field is unset
		if val == "" {
			continue
		}

		handler, ok := handlers[strings.ToLower(colLookup[row])]
		if !ok {
			fmt.Fprintf(os.Stderr, "no handler for field %d (%q); ignoring\n", row, colLookup[row])

			handler = unimplementedHandler

			// Suppress for next coin
			handlers[strings.ToLower(colLookup[row])] = handler
		}

		err := handler(coin, val)
		if err != nil {
			return err
		}
	}

	return nil
}

func recordID(coin *simplenuds.NUDS, val string) error {
	coin.Control.RecordID = val
	return nil
}

// Single image for entire coin.  For example,
// http://numismatics.org/collection/1922.999.73.xml shows
// <digRep>
//   <mets:fileSec>
//     <mets:fileGrp USE="obverse">
//        <mets:file USE="archive" MIMETYPE="image/jpeg">
//           <mets:FLocat LOCYPE="URL"
//              xlink:href="http://numismatics.org/collectionimages/19001949/1922/1922.999.73.obv.noscale.jpg"/>
//        </mets:file>
// ...
// mets is the namespace http://www.loc.gov/METS/
// and the mets schema seems to be http://www.loc.gov/standards/mets/mets.xsd
func coinSingleURLImageHandler(coin *simplenuds.NUDS, val string) error {
	// Create a file group if one doesn't exist
	if len(coin.DefaultDigRep().FileSec.FileGrp) == 0 {
		coin.DefaultDigRep().FileSec.FileGrp = []simplenuds.FileGrp{
			{
				File: []simplenuds.File{},
				USE:  "combined", // This handler is for single URLs, so "combined"
			},
		}
	}

	coin.DefaultDigRep().FileSec.FileGrp[0].AppendFile(
		simplenuds.File{
			USE: "reference",
			FLocat: []simplenuds.FLocat{
				{
					LOCTYPE: "URL",
					Href:    val,
				},
			},
		},
	)

	return nil
}

// For example,
// <nuds>
//   <descMeta>
//     <title xml:lang="en">Silver drahm of Khusraw II, MR, AD 591 - 628. 1922.999.73</title>
//     <subjectSet/>
//     <typeDesc>
//       <objectType xlink:href="http://nomisma.org/id/coin" xlink:type="simple">Coin</objectType>
//       <denomination>drahm</denomination>
func denominationHandler(coin *simplenuds.NUDS, val string) error {
	// TODO produce structured data for well-known types such as drachm
	coin.DescMeta.TypeDesc.AppendDenomination(simplenuds.Denomination(val))
	return nil
}

// For example,
// <nuds>
//   <descMeta>
//     <title xml:lang="en">Silver drahm of Khusraw II, MR, AD 591 - 628. 1922.999.73</title>
//     <subjectSet/>
//     <typeDesc>
//       <objectType xlink:href="http://nomisma.org/id/coin" xlink:type="simple">Coin</objectType>
//       <material xlink:href="http://nomisma.org/id/ar" xlink:type="simple">Silver</material>
func metalHandler(coin *simplenuds.NUDS, val string) error {
	material, err := getMaterial(val)
	if err != nil {
		return err
	}

	coin.DescMeta.TypeDesc.AppendMaterial(material)

	return nil
}

func diameterInMMHandler(coin *simplenuds.NUDS, val string) error {
	coin.DescMeta.DefaultPhysDesc().DefaultMeasurementsSet().Diameter = &simplenuds.Diameter{
		Units: "mm",
		Value: val,
	}

	return nil
}

func weightHandler(coin *simplenuds.NUDS, val string) error {
	// Rewrite European comma-separated such as "3,7"
	val = strings.Replace(val, ",", ".", 1)

	// Validate data
	if _, err := strconv.ParseFloat(val, 32); err != nil {
		fmt.Fprintf(os.Stderr, "invalid weight %q; ignoring\n", val)
	}

	coin.DescMeta.DefaultPhysDesc().DefaultMeasurementsSet().Weight = &simplenuds.Weight{
		Units: "g",
		Value: val,
	}

	return nil
}

func mintHandler(coin *simplenuds.NUDS, val string) error {
	// TODO.  nolint:godox
	// http://numismatics.org/collection/1960.10.1.xml
	// has
	// <geographic>
	//   <geogname xlink:role="region" xlink:type="simple">Mashriq</geogname>
	//   <geogname xlink:role="locality" xlink:type="simple">uncertain</geogname>
	// </geographic>
	return nil
}

func titleHandler(coin *simplenuds.NUDS, val string) error {
	coin.DescMeta.DefaultTitle()[0] = simplenuds.Title{
		Lang:  "en",
		Value: val,
	}

	return nil
}

// Example
// <control>
//   <rightsStmt>
//     <license for="data" xlink:type="simple"
//         xlink:href="http://opendatacommons.org/licenses/odbl/">Metadata are openly licensed with a
//              Open Data Commons Open Database License (ODbL)</license>
//     <license for="images" xlink:type="simple"
//         xlink:href="https://creativecommons.org/choose/mark/">Public Domain Mark</license>
//     <rights xlink:type="simple"
//         xlink:href="http://rightsstatements.org/vocab/NoC-US/1.0/">No Copyright - United States</rights>
// </rightsStmt>
func rightsURLHandler(coin *simplenuds.NUDS, val string) error {

	// warn and ignore if val is not an URL
	_, err := url.ParseRequestURI(val)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q is not a valid URL\n", val)
		return nil
	}

	// This implementation gives the same right to data and images
	coin.Control.RightsStmt.AppendLicense(simplenuds.License{
		For:  "data",
		Type: "simple",
		Href: val,
	})
	coin.Control.RightsStmt.AppendLicense(simplenuds.License{
		For:  "images",
		Type: "simple",
		Href: val,
	})

	return nil
}

// The source of the data, e.g. the agent
func sourceHandler(coin *simplenuds.NUDS, val string) error {
	// Note that `AgencyName` appears on the admin screen,
	// http://localhost:9080/orbeon/numishare/admin/edit/coin/?id=215654
	// but the screen users see, e.g.
	// http://localhost:9080/orbeon/numishare/collection1/id/215654
	// will not show it.
	// will not show it, although it will "export".
	coin.Control.MaintenanceAgency.AgencyName.Value = val
	return nil
}

// When this record was created digitally for the first time
func recordCreatedDateHandler(coin *simplenuds.NUDS, val string) error {
	// Note that `EventDateTime` appears on the admin screen,
	// http://localhost:9080/orbeon/numishare/admin/edit/coin/?id=215654
	// but the screen users see, e.g.
	// http://localhost:9080/orbeon/numishare/collection1/id/215654
	// will not show it, although it will "export".
	creationEvent := coin.Control.MaintenanceHistory.
		GetOrCreateEventType("created")
	creationEvent.EventDateTime.Value = val
	creationEvent.EventDateTime.StandardDateTime = val

	return nil
}

// Who created the original record
func reporterHandler(coin *simplenuds.NUDS, val string) error {
	// Note that `AgencyName` appears on the admin screen,
	// http://localhost:9080/orbeon/numishare/admin/edit/coin/?id=215654
	// but the screen users see, e.g.
	// http://localhost:9080/orbeon/numishare/collection1/id/215654
	// will not show it.
	// will not show it, although it will "export".
	creationEvent := coin.Control.MaintenanceHistory.
		GetOrCreateEventType("created")

	creationEvent.Agent.Value = val

	// We assume all Zeno.ru records are created by humans
	creationEvent.AgentType.Value = "human"

	return nil
}

func detailsHandler(coin *simplenuds.NUDS, val string) error {
	coin.DescMeta.AppendDescriptionSet(
		simplenuds.DescriptionSet{
			Description: []simplenuds.Description{
				{
					Value: val,
				},
			},
		},
	)

	return nil
}

func unimplementedHandler(coin *simplenuds.NUDS, val string) error {
	return nil
}

func getMaterial(material string) (simplenuds.Material, error) {
	HRefs := map[string]string{
		"AR": "http://nomisma.org/id/ar",
		"AV": "http://nomisma.org/id/av",
		// TODO Structured types for other common metals nolint:godox
	}
	Texts := map[string]string{
		"AR": "Silver",
		"AV": "Gold",
	}

	HRef, ok := HRefs[material]
	if !ok {
		// Warning
		fmt.Fprintf(os.Stderr, "unimplemented metal: %q\n", material)

		return simplenuds.Material{
			Text: material,
		}, nil
	}

	Text, ok := Texts[material]
	if !ok {
		return simplenuds.Material{}, fmt.Errorf("unimplemented metal: %q", material)
	}

	return simplenuds.Material{
		HRef: HRef,
		Type: "simple",
		Text: Text,
	}, nil
}
