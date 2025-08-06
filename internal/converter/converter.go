package converter

import (
	"fmt"

	"opcmss/internal/model"
)

// ConvertModbusToOPC converts a Modbus tag to an OPC UA tag with configurable prefix
func ConvertModbusToOPC(modbusTag model.ModbusTag, namespaceIndex uint16, prefix string) model.OPCTag {
	// Generate OPC UA NodeID with the configurable prefix
	// Format: ns=<namespace>;s=<prefix><tag_name>
	nodeID := fmt.Sprintf("ns=%d;s=%s%s", namespaceIndex, prefix, modbusTag.Name)

	// Determine data type based on register type and size
	dataType := determineDataType(modbusTag.RegisterType, modbusTag.Size)

	return model.OPCTag{
		Name:     modbusTag.Name,
		NodeID:   nodeID,
		DataType: dataType,
	}
}

// determineDataType maps Modbus register types to OPC data types
func determineDataType(registerType string, size uint16) string {
	switch registerType {
	case "HoldingRegister", "InputRegister":
		switch size {
		case 1:
			return "INT" // 16-bit integer
		case 2:
			return "REAL" // 32-bit float (2 registers)
		}
		return "INT"
	case "Coil", "DiscreteInput":
		return "BOOL"
	default:
		return "INT"
	}
}

// ConvertAllModbusToOPC converts all Modbus tags to OPC tags with configurable prefix
func ConvertAllModbusToOPC(modbusTags []model.ModbusTag, namespaceIndex uint16, prefix string) []model.OPCTag {
	opcTags := make([]model.OPCTag, len(modbusTags))

	for i, modbusTag := range modbusTags {
		opcTags[i] = ConvertModbusToOPC(modbusTag, namespaceIndex, prefix)
	}

	return opcTags
}
