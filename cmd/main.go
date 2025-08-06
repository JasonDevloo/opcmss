package main

import (
	"fmt"
	"log"

	"opcmss/internal/converter"
	"opcmss/internal/modbus"
	"opcmss/internal/opcua"
	"opcmss/internal/parser"
)

const (
	// OPC UA configuration constants
	OPC_NAMESPACE_INDEX = 4
	OPC_NODE_PREFIX     = "|var|NEXTO PLC.Z.O83."
	OPC_ENDPOINT        = "opc.tcp://172.29.48.69:4840"

	MODBUS_ENDPOINT  = "172.29.48.69:502"
	MODBUS_TAGS_FILE = "/home/maimus/GoProjects/OPCvsModSymSrv/cmd/example_tags.tsv"

	// Comparison settings
	TAGS_TO_COMPARE = 20
)

func main() {
	// Parse Modbus tags from TSV using the proper TSV parser
	modbusTags, err := parser.ParseTagsTSV(MODBUS_TAGS_FILE)
	if err != nil {
		log.Fatal("Failed to parse TSV:", err)
	}

	// Convert Modbus tags to OPC tags using constants
	opcTags := converter.ConvertAllModbusToOPC(modbusTags, OPC_NAMESPACE_INDEX, OPC_NODE_PREFIX)

	// Create OPC UA client
	client, err := opcua.NewClient(OPC_ENDPOINT)
	if err != nil {
		log.Fatal("Failed to create OPC client:", err)
	}
	defer client.Close()

	// Create Modbus client
	modbusClient, err := modbus.NewClient(MODBUS_ENDPOINT)
	if err != nil {
		log.Fatal("Failed to create Modbus client:", err)
	}
	defer modbusClient.Close()

	totalTags := len(modbusTags)
	fmt.Printf("Total tags available: %d\n", totalTags)
	fmt.Printf("Comparing %d evenly spaced tags:\n\n", TAGS_TO_COMPARE)

	// Calculate step size for even spacing
	step := totalTags / TAGS_TO_COMPARE
	if step < 1 {
		step = 1
	}

	successCount := 0
	errorCount := 0

	// Compare evenly spaced tags
	for i := 0; i < TAGS_TO_COMPARE && i*step < totalTags; i++ {
		index := i * step

		// Ensure we don't exceed array bounds
		if index >= totalTags {
			break
		}

		opcTag := opcTags[index]
		modbusTag := modbusTags[index]

		fmt.Printf("=== Tag %d/%d (Index: %d) ===\n", i+1, TAGS_TO_COMPARE, index)
		fmt.Printf("Name: %s\n", opcTag.Name)
		fmt.Printf("Type: %s, Address: %d, Size: %d\n", modbusTag.RegisterType, modbusTag.Address, modbusTag.Size)

		// Read from OPC UA
		fmt.Printf("OPC UA: ")
		opcValue, opcErr := client.ReadTag(opcTag)
		if opcErr != nil {
			fmt.Printf("Error: %v\n", opcErr)
			errorCount++
		} else {
			fmt.Printf("Value: %v (type: %T)\n", opcValue, opcValue)
			successCount++
		}

		// Read from Modbus
		fmt.Printf("Modbus: ")
		modbusValue, modbusErr := modbusClient.ReadTag(modbusTag)
		if modbusErr != nil {
			fmt.Printf("Error: %v\n", modbusErr)
		} else {
			fmt.Printf("Value: %v\n", modbusClient.FormatTagValue(modbusTag, modbusValue))
		}

		// Compare values if both successful
		if opcErr == nil && modbusErr == nil {
			if compareValues(opcValue, modbusValue, modbusTag.RegisterType) {
				fmt.Printf("✓ Values match!\n")
			} else {
				fmt.Printf("✗ Values differ!\n")
			}
		}

		fmt.Println()
	}

	fmt.Printf("Summary: %d successful OPC reads, %d errors out of %d attempts\n",
		successCount, errorCount, TAGS_TO_COMPARE)
}

// compareValues compares OPC and Modbus values based on the register type
func compareValues(opcValue, modbusValue any, registerType string) bool {
	switch registerType {
	case "Coil", "DiscreteInput":
		opcBool, opcOk := opcValue.(bool)
		modbusBool, modbusOk := modbusValue.(bool)
		return opcOk && modbusOk && opcBool == modbusBool

	case "HoldingRegister", "InputRegister":
		// Handle different numeric types
		opcFloat := convertToFloat64(opcValue)
		modbusFloat := convertToFloat64(modbusValue)

		if opcFloat == nil || modbusFloat == nil {
			return false
		}

		// Allow small floating point differences
		diff := *opcFloat - *modbusFloat
		if diff < 0 {
			diff = -diff
		}
		return diff < 0.001

	default:
		return fmt.Sprintf("%v", opcValue) == fmt.Sprintf("%v", modbusValue)
	}
}

// convertToFloat64 converts various numeric types to float64
func convertToFloat64(value any) *float64 {
	switch v := value.(type) {
	case float32:
		f := float64(v)
		return &f
	case float64:
		return &v
	case int16:
		f := float64(v)
		return &f
	case int32:
		f := float64(v)
		return &f
	case int:
		f := float64(v)
		return &f
	case uint16:
		f := float64(v)
		return &f
	case uint32:
		f := float64(v)
		return &f
	default:
		return nil
	}
}
