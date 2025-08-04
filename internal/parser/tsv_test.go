package parser

import (
	"os"
	"testing"
)

func TestLoadTSV(t *testing.T) {
	tsv := `Name	Address	DataType	Description
Temperature	40001	REAL	Main Temp
PumpStatus	40003	BOOL	Pump On/Off`

	tmp := "test_tags.tsv"
	err := os.WriteFile(tmp, []byte(tsv), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)

	tags, err := ParseTagsTSV(tmp)
	if err != nil {
		t.Fatalf("Error loading TSV: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("Expected 2 tags, got %d", len(tags))
	}

	if tags[0].Name != "Temperature" || tags[1].RegisterType != "BOOL" {
		t.Errorf("Parsed values incorrect: %+v", tags)
	}
}
