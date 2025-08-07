package parser

import (
	"os"
	"strings"
	"testing"

	"opcmss/internal/model"
)

func TestParseTagsTSV_ValidData(t *testing.T) {
	// TSV with correct format: Name, RegisterType, Address, ModbusAddress, Size, Range
	tsv := `Temperature	HoldingRegister	1	40001	2	1..2
PumpStatus	Coil	5	00005	1	5..5
Pressure	HoldingRegister	10	40010	1	10..10`

	tmp := "test_valid_tags.tsv"
	err := os.WriteFile(tmp, []byte(tsv), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)

	tags, err := ParseTagsTSV(tmp)
	if err != nil {
		t.Fatalf("Error loading TSV: %v", err)
	}

	if len(tags) != 3 {
		t.Fatalf("Expected 3 tags, got %d", len(tags))
	}

	// Verify first tag
	expectedTag1 := model.ModbusTag{
		Name:          "Temperature",
		RegisterType:  "HoldingRegister",
		Address:       1,
		ModbusAddress: 40001,
		Size:          2,
		Range:         "1..2",
	}
	if tags[0] != expectedTag1 {
		t.Errorf("First tag incorrect. Expected: %+v, Got: %+v", expectedTag1, tags[0])
	}

	// Verify second tag
	expectedTag2 := model.ModbusTag{
		Name:          "PumpStatus",
		RegisterType:  "Coil",
		Address:       5,
		ModbusAddress: 5,
		Size:          1,
		Range:         "5..5",
	}
	if tags[1] != expectedTag2 {
		t.Errorf("Second tag incorrect. Expected: %+v, Got: %+v", expectedTag2, tags[1])
	}
}

func TestParseTagsTSV_EmptyFile(t *testing.T) {
	tmp := "test_empty.tsv"
	err := os.WriteFile(tmp, []byte(""), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)

	tags, err := ParseTagsTSV(tmp)
	if err != nil {
		t.Fatalf("Expected no error for empty file, got: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("Expected 0 tags for empty file, got %d", len(tags))
	}
}

func TestParseTagsTSV_InconsistentFieldCount(t *testing.T) {
	// TSV with inconsistent field counts should return an error
	tsv := `ValidTag	HoldingRegister	1	40001	1	1..1
IncompleteRecord	HoldingRegister
AnotherValidTag	Coil	2	00002	1	2..2`

	tmp := "test_inconsistent.tsv"
	err := os.WriteFile(tmp, []byte(tsv), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)

	_, err = ParseTagsTSV(tmp)
	if err == nil {
		t.Fatal("Expected error for inconsistent field count, got none")
	}

	// Should get a CSV parsing error about wrong number of fields
	if !strings.Contains(err.Error(), "wrong number of fields") {
		t.Errorf("Expected 'wrong number of fields' error, got: %v", err)
	}
}

func TestParseTagsTSV_InvalidNumbers(t *testing.T) {
	// Records with invalid numeric values should be skipped
	tsv := `GoodTag	HoldingRegister	1	40001	1	1..1
BadAddress	HoldingRegister	NOTNUMBER	40002	1	2..2
BadModbusAddr	HoldingRegister	3	INVALID	1	3..3
BadSize	HoldingRegister	4	40004	BAD	4..4
AnotherGoodTag	Coil	5	00005	1	5..5`

	tmp := "test_invalid_numbers.tsv"
	err := os.WriteFile(tmp, []byte(tsv), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)

	tags, err := ParseTagsTSV(tmp)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should only parse the 2 valid records, skip ones with invalid numbers
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags (skipping invalid numbers), got %d", len(tags))
	}

	if len(tags) >= 2 {
		if tags[0].Name != "GoodTag" || tags[1].Name != "AnotherGoodTag" {
			t.Errorf("Invalid number records not properly skipped. Got tags: %v", tags)
		}
	}
}

func TestParseTagsTSV_FileNotFound(t *testing.T) {
	// Test with non-existent file
	_, err := ParseTagsTSV("nonexistent_file.tsv")
	if err == nil {
		t.Fatal("Expected error for non-existent file, got none")
	}

	// Error should be related to file opening
	if !os.IsNotExist(err) {
		t.Errorf("Expected file not found error, got: %v", err)
	}
}

func TestParseTagsTSV_WhitespaceHandling(t *testing.T) {
	// Test with various whitespace around values (but exactly 6 fields per record)
	// Creating TSV with leading/trailing spaces but consistent tab separators
	lines := []string{
		"  SpacedName  \t  HoldingRegister  \t1\t40001\t1\t  1..1  ",
		" TabName \tCoil\t2\t00002\t1\t2..2",
		" MixedSpaces \t HoldingRegister \t3\t40003\t2\t 3..4 ",
	}
	tsv := strings.Join(lines, "\n")

	tmp := "test_whitespace.tsv"
	err := os.WriteFile(tmp, []byte(tsv), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)

	tags, err := ParseTagsTSV(tmp)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(tags) != 3 {
		t.Fatalf("Expected 3 tags, got %d", len(tags))
	}

	// Verify whitespace is properly trimmed
	expectedNames := []string{"SpacedName", "TabName", "MixedSpaces"}
	expectedRegTypes := []string{"HoldingRegister", "Coil", "HoldingRegister"}
	expectedRanges := []string{"1..1", "2..2", "3..4"}

	for i, tag := range tags {
		if tag.Name != expectedNames[i] {
			t.Errorf("Tag %d name: expected '%s', got '%s'", i, expectedNames[i], tag.Name)
		}
		if tag.RegisterType != expectedRegTypes[i] {
			t.Errorf("Tag %d register type: expected '%s', got '%s'", i, expectedRegTypes[i], tag.RegisterType)
		}
		if tag.Range != expectedRanges[i] {
			t.Errorf("Tag %d range: expected '%s', got '%s'", i, expectedRanges[i], tag.Range)
		}
	}
}

func TestParseTagsTSV_NumericBoundaries(t *testing.T) {
	// Test with boundary values for uint16 and uint32
	tsv := `MinValues	HoldingRegister	0	0	0	0..0
MaxUint16	HoldingRegister	65535	4294967295	65535	1..65535
MaxValidValues	Coil	32767	2147483647	32767	1..32767`

	tmp := "test_boundaries.tsv"
	err := os.WriteFile(tmp, []byte(tsv), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)

	tags, err := ParseTagsTSV(tmp)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(tags) != 3 {
		t.Fatalf("Expected 3 tags, got %d", len(tags))
	}

	// Check boundary values
	if tags[0].Address != 0 || tags[0].ModbusAddress != 0 || tags[0].Size != 0 {
		t.Errorf("Min values not parsed correctly: %+v", tags[0])
	}

	if tags[1].Address != 65535 || tags[1].ModbusAddress != 4294967295 || tags[1].Size != 65535 {
		t.Errorf("Max values not parsed correctly: %+v", tags[1])
	}
}

func TestParseTagsTSV_NumberOverflow(t *testing.T) {
	// Test with numbers that exceed uint16/uint32 limits - should be skipped
	tsv := `ValidTag	HoldingRegister	1	40001	1	1..1
OverflowAddress	HoldingRegister	70000	40002	1	2..2
OverflowModbus	HoldingRegister	3	5000000000	1	3..3
OverflowSize	HoldingRegister	4	40004	70000	4..4
AnotherValidTag	Coil	5	00005	1	5..5`

	tmp := "test_overflow.tsv"
	err := os.WriteFile(tmp, []byte(tsv), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)

	tags, err := ParseTagsTSV(tmp)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should only parse the 2 valid records, skip overflow values
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags (skipping overflow values), got %d", len(tags))
	}

	if len(tags) >= 2 {
		if tags[0].Name != "ValidTag" || tags[1].Name != "AnotherValidTag" {
			t.Errorf("Overflow records not properly skipped. Got tags: %v", tags)
		}
	}
}
