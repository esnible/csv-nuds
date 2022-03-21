package converter

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name   string
		coin   map[string]string
		golden string
	}{
		{
			name: "zeno 264199",
			coin: map[string]string{
				"creationtime": "13 Dec 20 11:55:36 +0300",
				"denomination": "Drakhm",
				"diameter":     "29",
				"id":           "264199",
				"imageurl":     "https://zeno.ru/data/2807/medium/Kaykhusru-24.jpg",
				"metal":        "AR",
				"mint":         "?",
				"reporter":     "Ombo",
				"rightsurl":    "https://rightsstatements.org/page/CNE/1.0/?language=en",
				"source":       "Zeno.ru",
				"title":        "Sasanid , Kaykhusru 2 , AR drakhme",
				"weight":       "3.62",
			},
			golden: "nuds264199.xml",
		},
	}

	converter := NewConverter(time.Time{})
	for n, testcase := range tests {
		nuds, err := converter.GenerateNUDS(testcase.coin)
		if err != nil {
			t.Fatalf("Failed converting test %d: %v", n, err)
		}

		data, err := xml.MarshalIndent(nuds, " ", "  ")
		if err != nil {
			log.Fatal(err)
		}
		got := string(data)

		want := goldenValue(t, testcase.golden, got, *update)
		if got != want {
			t.Errorf("Want:\n%s\nGot:\n%s", want, got)
		}
	}
}

func goldenValue(t *testing.T, goldenFile string, actual string, update bool) string {
	t.Helper()
	goldenPath := "testdata/" + goldenFile + ".golden"

	f, err := os.OpenFile(goldenPath, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", goldenPath, err)
	}
	defer f.Close()

	if update {
		_, err := f.WriteString(actual)
		if err != nil {
			t.Fatalf("Error writing to file %s: %s", goldenPath, err)
		}

		return actual
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Error reading file %s: %s", goldenPath, err)
	}
	return string(content)
}
