
# csv-nuds

cvs-nuds is a proof-of-concept for a workflow to convert from numismatic data in spreadsheets to [NUDS](https://nomisma.org/).

Once the data is in NUDS format it can be loaded into [Numishare](https://github.com/ewg118/numishare) or other software that handles the NUDS format.

**Currently it isn't working.**

Currently a lot of data is numismatic data is stored in relational databases or spreadsheets.  Many of these formats make it easy to export to flat [CSV](https://en.wikipedia.org/wiki/Comma-separated_values) files.  Those files are easy to manipulate manually or with scripts.

The idea behind this tool is to produce XML that is valid to [the NUDS schema](http://nomisma.org/nuds.xsd) so that it can be loaded into Numishare or as an interchange format between different systems.

## Implementation

I was unable to generate language bindings from http://nomisma.org/nuds.xsd so I have created a simple subset of NUDS by hand.

I have custom code to move data from named columns in CSV files into NUDS.

### Tutorial / Regenerating test data

Install [Go](https://en.wikipedia.org/wiki/Go_(programming_language)

Execute `go run csv2nuds.go data/zeno.csv data/every-zeno.csv > data/zeno.nuds` to convert 20 records from an ad-hoc CSV file into NUDS XML.

Note: This data was manually scraped from [https://zeno.ru/](https://zeno.ru/).  It's just 20 random Khusru II drachms.  If anyone has public-domain or Creative Commons numismatic data in CSV format please let me know.

In addition to writing the data the tool currently outputs

```
no handler for field 1 ("url"); ignoring
no handler for field 3 ("date"); ignoring
no handler for field 10 ("reporter"); ignoring
no handler for field 11 ("reporterUrl"); ignoring
no handler for field 0 ("source"); ignoring
no handler for field 1 ("rights"); ignoring
no handler for field 12 ("additionalDetails"); ignoring
unimplemented metal: "silver washed AE"
unimplemented metal: "Tin-zinc alloy"
```

This is because I couldn't quite figure out the best way to express those concepts in NUDS.  The next version should support all of the demo fields.

An example Sasanian drachm in [the ANS collection](http://numismatics.org/search/) can be fetched from their server.

`curl http://numismatics.org/collection/1922.999.73.xml > data/1922.999.73.xml`

## Applying NUDS to a Numishare server

I have not been able to successfully upload the data.  The plan is that the generated NUDS will be easy to load upon Numishare with a script like this:

```
EXIST_HOST=localhost:8888
COLLECTION=collection1
EXIST_USER=admin
EXIST_PASSWORD=
curl -v -X POST --user "admin:" http://"$EXIST_HOST"/exist/rest/db/"$COLLECTION"/objects/ --upload-file data/zeno.nuds
```

I also failed to upload interactively.  According to http://$EXIST_HOST:8888/exist/apps/doc/uploading-files "eXist-db's Dashboard's comes with a Collections pane."  (I don't have a Collections pane.)  That page also suggests using eXist-db's built-in Integrated Development Environment (IDE). File, Manage from the menu; click on the Upload button.  This sequence of steps produces no error but also nothing new in Numishare.