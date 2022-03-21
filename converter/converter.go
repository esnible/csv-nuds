// Convert columnar data to NUDS

package converter

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/esnible/csv-nuds/simplenuds"
)

type NUDSWriter func(coin *simplenuds.NUDS, val string) error

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

type Converter struct {
	// Handlers for the different column names
	Handlers map[string]NUDSWriter

	// The time to use when creating records
	Timestamp time.Time
}

func NewConverter(timestamp time.Time) Converter {
	return Converter{

		Handlers: map[string]NUDSWriter{
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
		},
		Timestamp: timestamp,
	}
}

// GenerateNUDS() generates NUDS from a slice of column values (a CSV coin row) and optional second row
func (converter *Converter) GenerateNUDS(coin map[string]string) (*simplenuds.NUDS, error) {
	retval := simplenuds.NewNUDS("physical", converter.Timestamp)

	for key, val := range coin {
		handler, ok := converter.Handlers[key]
		if !ok {
			fmt.Fprintf(os.Stderr, "no handler for %q, a %q; ignoring\n", val, key)

			handler = unimplementedHandler

			// Suppress for next coin
			converter.Handlers[key] = handler
		}

		err := handler(&retval, val)
		if err != nil {
			return nil, err
		}
	}

	return &retval, nil
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
