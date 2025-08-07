package parser

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"opcmss/internal/model"
)

func ParseTagsTSV(filename string) ([]model.ModbusTag, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // Tab-separated values

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var tags []model.ModbusTag
	for _, record := range records {
		if len(record) < 1 {
			continue // Skip empty records
		}

		address, err := strconv.ParseUint(record[2], 10, 16)
		if err != nil {
			continue
		}

		modbusAddress, err := strconv.ParseUint(record[3], 10, 32)
		if err != nil {
			continue
		}

		size, err := strconv.ParseUint(record[4], 10, 16)
		if err != nil {
			continue
		}

		tag := model.ModbusTag{
			Name:          strings.TrimSpace(record[0]),
			RegisterType:  strings.TrimSpace(record[1]),
			Address:       uint16(address),
			ModbusAddress: uint32(modbusAddress),
			Size:          uint16(size),
			Range:         strings.TrimSpace(record[5]),
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
