
# csv-nuds

_csv-nuds_ is a proof-of-concept for a workflow to convert from numismatic data in spreadsheets to [NUDS](http://www.greekcoinage.org/nuds.html), a specification from [Nomisma.org](http://nomisma.org/).  (Pronunced "N.U.D.S.").

Once the data is in NUDS format it can be loaded into [Numishare](https://github.com/ewg118/numishare) or other software that handles the NUDS format.

A lot of numismatic data is stored in relational databases or spreadsheets.  Many of these formats make it easy to export to flat [CSV](https://en.wikipedia.org/wiki/Comma-separated_values) files.  Those files are easy to manipulate manually or with scripts.

The idea behind csv-nuds is to export all the data currently in flat files and databases and produce XML that is valid to [the NUDS schema](http://nomisma.org/nuds.xsd).  NUDS data can be loaded into Numishare or as an interchange format between different systems.

### Tutorial / Regenerating NUDS from tutorial sample data

Install [Go](https://en.wikipedia.org/wiki/Go_(programming_language))

Execute `go run csv2nuds.go zeno data/zeno.csv data/every-zeno.csv` to convert 20 records from an ad-hoc CSV file into 20 NUDS XML files.

Note: This data was manually scraped from [https://zeno.ru/](https://zeno.ru/).  It's just 20 random Khusru II drachms.  If anyone has public-domain or Creative Commons numismatic data in CSV format please let me know.

## Applying NUDS to a Numishare server

Generated NUDS can be be sent to Numishare with a script like this:

```
EXIST_HOST=localhost:8888
COLLECTION=collection1
EXIST_USER=admin
EXIST_PASSWORD=
for filename in zeno/*.xml; do
   curl -v --user "$EXIST_USER":"$EXIST_PASSWORD" http://"$EXIST_HOST"/exist/rest/db/"$COLLECTION"/objects/ --upload-file "$filename"
done 
```

## Limitations of the convert / questions about how to represent

In addition to writing the data the tool currently outputs

```
no handler for field 1 ("url"); ignoring
no handler for field 3 ("date"); ignoring
no handler for field 11 ("reporterUrl"); ignoring
unimplemented metal: "silver washed AE"
unimplemented metal: "Tin-zinc alloy"
```

- I'll like make the Zeno `reporterUrl` into an `<acknowledgement>`.  (I originally considered `<copyrightHolder>` (even though it might not be), or perhaps `<owner>`).  None of these appear in Numishare (at this time.)
- I am not sure what to make the Zeno `url` into.  Zeno itself might be a `<collection>` (but of images, not coins).  There should be some kind of way to refer/link to another representation of the same object, but I don't know it.
- The Zeno category (not currently in the CSV) will become a `<department>`.  There will be thousands of them.
- Numishare has `<material xlink:href="http://nomisma.org/id/sn" xlink:type="simple">Tin</material>` but nothing for a Tin-zinc alloy.
- I don't know how to represent "silver washed AE"

For comparison between this tool's output and "real NUDS", an example Sasanian drachm in [the ANS collection](http://numismatics.org/search/) can be fetched from their server.

`curl http://numismatics.org/collection/1922.999.73.xml > data/1922.999.73.xml`

My goal is for this tool to produce XML with a similar level of complexity.

## Implementation

I was unable to generate language bindings from http://nomisma.org/nuds.xsd so I have created a simple subset of NUDS by hand.

I have custom code to move data from named columns in CSV files into NUDS.

### Testing

We test the code with `go test -v ./...`

The test compares a coin expressed in key/value pairs with NUDS XML stored in a golden file.  Expected NUDS XML are stored in the [converter/testdata](converter/testdata) folder.
