package modbus

import (
	"testing"

	"opcmss/internal/model"
)

func TestModbusConnection(t *testing.T) {
	client, err := NewClient("172.29.48.69:502")
	if err != nil {
		t.Fatalf("Modbus server not available: %v", err)
	}
	defer client.Close()

	tag := model.ModbusTag{
		Name:          "TestHR",
		RegisterType:  "HoldingRegister",
		Address:       1,
		ModbusAddress: 400001,
		Size:          1,
		Range:         "1..1",
	}

	val, err := client.ReadTag(tag)
	if err != nil {
		t.Errorf("Failed to read tag: %v", err)
	}

	t.Logf("Read value: %v", val)
}
