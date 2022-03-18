
# csv-nuds

_csv-nuds_ is a proof-of-concept for a workflow to convert from numismatic data in spreadsheets to [NUDS](http://www.greekcoinage.org/nuds.html), a specification from [Nomisma.org](http://nomisma.org/)

Once the data is in NUDS format it can be loaded into [Numishare](https://github.com/ewg118/numishare) or other software that handles the NUDS format.

**Currently it isn't working.**

A lot of numismatic data is stored in relational databases or spreadsheets.  Many of these formats make it easy to export to flat [CSV](https://en.wikipedia.org/wiki/Comma-separated_values) files.  Those files are easy to manipulate manually or with scripts.

The idea behind csv-nuds is to export all the data currently in flat files and databases and produce XML that is valid to [the NUDS schema](http://nomisma.org/nuds.xsd).  NUDS data can be loaded into Numishare or as an interchange format between different systems.

## Implementation

I was unable to generate language bindings from http://nomisma.org/nuds.xsd so I have created a simple subset of NUDS by hand.

I have custom code to move data from named columns in CSV files into NUDS.

### Tutorial / Regenerating test data

Install [Go](https://en.wikipedia.org/wiki/Go_(programming_language))

Execute `go run csv2nuds.go data/zeno.csv data/every-zeno.csv > data/zeno.nuds` to convert 20 records from an ad-hoc CSV file into NUDS XML.
**Note that this XML is not a valid Document, as it contains more than one
root.  It is a valid 'fragment'.**

Note: This data was manually scraped from [https://zeno.ru/](https://zeno.ru/).  It's just 20 random Khusru II drachms.  If anyone has public-domain or Creative Commons numismatic data in CSV format please let me know.

In addition to writing the data the tool currently outputs

```
no handler for field 1 ("url"); ignoring
no handler for field 3 ("date"); ignoring
no handler for field 10 ("reporter"); ignoring
no handler for field 11 ("reporterUrl"); ignoring
no handler for field 0 ("source"); ignoring
no handler for field 12 ("additionalDetails"); ignoring
unimplemented metal: "silver washed AE"
unimplemented metal: "Tin-zinc alloy"
```

This is because I couldn't quite figure out the best way to express those concepts in NUDS.  The next version should support all of the demo fields.

Execute `go run csv2nuds.go data/58627.csv data/every-zeno.csv > data/58627.xml`

For comparison between this tool's output and "real NUDS", an example Sasanian drachm in [the ANS collection](http://numismatics.org/search/) can be fetched from their server.

`curl http://numismatics.org/collection/1922.999.73.xml > data/1922.999.73.xml`

My goal is for this tool to produce XML with a similar level of complexity.

## Applying NUDS to a Numishare server

I have not been able to successfully upload generated NUDS data to Numishare.
The plan is that the generated NUDS will be easy to load upon Numishare with a script like this:

```
EXIST_HOST=localhost:8888
COLLECTION=collection1
EXIST_USER=admin
EXIST_PASSWORD=
curl -v -X POST --user "admin:" http://"$EXIST_HOST"/exist/rest/db/"$COLLECTION"/objects/ --upload-file data/zeno.nuds
```
