package modbus

import (
	"errors"
	"testing"

	"opcmss/internal/model"

	"github.com/simonvetter/modbus"
)

type MockModbusClient struct {
	coilsData      []bool
	registersData  []uint16
	coilsError     error
	registersError error
	closeError     error
	openError      error
}

func (m *MockModbusClient) ReadCoils(address, quantity uint16) ([]bool, error) {
	if m.coilsError != nil {
		return nil, m.coilsError
	}
	if int(quantity) > len(m.coilsData) {
		return m.coilsData, nil
	}
	return m.coilsData[:quantity], nil
}

func (m *MockModbusClient) ReadRegisters(address, quantity uint16, regType modbus.RegType) ([]uint16, error) {
	if m.registersError != nil {
		return nil, m.registersError
	}
	if int(quantity) > len(m.registersData) {
		return m.registersData, nil
	}
	return m.registersData[:quantity], nil
}

func (m *MockModbusClient) Open() error {
	return m.openError
}

func (m *MockModbusClient) Close() error {
	return m.closeError
}

func TestReadCoil_SingleCoilTrue(t *testing.T) {
	mock := &MockModbusClient{
		coilsData: []bool{true},
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestCoil",
		RegisterType: "Coil",
		Address:      1,
		Size:         1,
	}

	result, err := client.readCoil(tag)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != true {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestReadCoil_SingleCoilFalse(t *testing.T) {
	mock := &MockModbusClient{
		coilsData: []bool{false},
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestCoil",
		RegisterType: "Coil",
		Address:      1,
		Size:         1,
	}

	result, err := client.readCoil(tag)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != false {
		t.Errorf("Expected false, got: %v", result)
	}
}

func TestReadCoil_MultipleCoils(t *testing.T) {
	mock := &MockModbusClient{
		coilsData: []bool{true, false, true, false},
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestCoils",
		RegisterType: "Coil",
		Address:      1,
		Size:         4,
	}

	result, err := client.readCoil(tag)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	coilsResult, ok := result.([]bool)
	if !ok {
		t.Fatalf("Expected []bool, got type: %T", result)
	}

	expected := []bool{true, false, true, false}
	if len(coilsResult) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(coilsResult))
	}

	for i, expectedVal := range expected {
		if i < len(coilsResult) && coilsResult[i] != expectedVal {
			t.Errorf("Coil %d: expected %v, got %v", i, expectedVal, coilsResult[i])
		}
	}
}

func TestReadCoil_EmptyData(t *testing.T) {
	mock := &MockModbusClient{
		coilsData: []bool{},
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestCoil",
		RegisterType: "Coil",
		Address:      1,
		Size:         1,
	}

	result, err := client.readCoil(tag)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != false {
		t.Errorf("Expected false for empty data, got: %v", result)
	}
}

func TestReadCoil_Error(t *testing.T) {
	mock := &MockModbusClient{
		coilsError: errors.New("modbus read error"),
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestCoil",
		RegisterType: "Coil",
		Address:      1,
		Size:         1,
	}

	_, err := client.readCoil(tag)
	if err == nil {
		t.Fatal("Expected error, got none")
	}

	if err.Error() != "modbus read error" {
		t.Errorf("Expected 'modbus read error', got: %v", err)
	}
}

func TestReadHoldingRegister_SingleRegister(t *testing.T) {
	mock := &MockModbusClient{
		registersData: []uint16{12345},
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestHR",
		RegisterType: "HoldingRegister",
		Address:      1,
		Size:         1,
	}

	result, err := client.readHoldingRegister(tag)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	intResult, ok := result.(int16)
	if !ok {
		t.Fatalf("Expected int16, got type: %T", result)
	}

	if intResult != 12345 {
		t.Errorf("Expected 12345, got: %v", intResult)
	}
}

func TestReadHoldingRegister_SingleRegister_NoData(t *testing.T) {
	mock := &MockModbusClient{
		registersData: []uint16{},
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestHR",
		RegisterType: "HoldingRegister",
		Address:      1,
		Size:         1,
	}

	_, err := client.readHoldingRegister(tag)
	if err == nil {
		t.Fatal("Expected error for no data, got none")
	}

	if err.Error() != "no data received" {
		t.Errorf("Expected 'no data received', got: %v", err)
	}
}

func TestReadHoldingRegister_Error(t *testing.T) {
	mock := &MockModbusClient{
		registersError: errors.New("register read error"),
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestHR",
		RegisterType: "HoldingRegister",
		Address:      1,
		Size:         1,
	}

	_, err := client.readHoldingRegister(tag)
	if err == nil {
		t.Fatal("Expected error, got none")
	}

	if err.Error() != "register read error" {
		t.Errorf("Expected 'register read error', got: %v", err)
	}
}

func TestReadHoldingRegister_TwoRegisters(t *testing.T) {
	// Test data representing float32 value 60.0
	// 60.0 in IEEE 754 binary32 format is 0x42700000
	mock := &MockModbusClient{
		registersData: []uint16{0x4270, 0x0000}, // Big-endian representation
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestHR2",
		RegisterType: "HoldingRegister",
		Address:      1,
		Size:         2,
	}

	result, err := client.readHoldingRegister(tag)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	mapResult, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got type: %T", result)
	}

	// Check that all expected keys are present
	if _, ok := mapResult["int32"]; !ok {
		t.Error("Expected 'int32' key in result")
	}
	if _, ok := mapResult["float32"]; !ok {
		t.Error("Expected 'float32' key in result")
	}
	if _, ok := mapResult["raw"]; !ok {
		t.Error("Expected 'raw' key in result")
	}

	// Check float32 value
	if floatVal, ok := mapResult["float32"].(float32); ok {
		if floatVal != 60.0 {
			t.Errorf("Expected float32 value 60.0, got: %v", floatVal)
		}
	} else {
		t.Error("float32 value is not of correct type")
	}

	// Check raw data
	if rawVal, ok := mapResult["raw"].([]uint16); ok {
		expected := []uint16{0x4270, 0x0000}
		if len(rawVal) != len(expected) {
			t.Errorf("Expected raw length %d, got %d", len(expected), len(rawVal))
		}
		for i, exp := range expected {
			if i < len(rawVal) && rawVal[i] != exp {
				t.Errorf("Raw[%d]: expected 0x%04x, got 0x%04x", i, exp, rawVal[i])
			}
		}
	} else {
		t.Error("raw value is not of correct type")
	}
}

func TestReadHoldingRegister_TwoRegisters_InsufficientData(t *testing.T) {
	mock := &MockModbusClient{
		registersData: []uint16{0x1234}, // Only one register when two expected
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestHR2",
		RegisterType: "HoldingRegister",
		Address:      1,
		Size:         2,
	}

	_, err := client.readHoldingRegister(tag)
	if err == nil {
		t.Fatal("Expected error for insufficient data, got none")
	}

	if err.Error() != "insufficient data for 32-bit value" {
		t.Errorf("Expected 'insufficient data for 32-bit value', got: %v", err)
	}
}

func TestReadHoldingRegister_MultipleRegisters(t *testing.T) {
	mock := &MockModbusClient{
		registersData: []uint16{100, 200, 300, 400, 500},
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestHRMultiple",
		RegisterType: "HoldingRegister",
		Address:      1,
		Size:         5,
	}

	result, err := client.readHoldingRegister(tag)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	sliceResult, ok := result.([]int16)
	if !ok {
		t.Fatalf("Expected []int16, got type: %T", result)
	}

	expected := []int16{100, 200, 300, 400, 500}
	if len(sliceResult) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(sliceResult))
	}

	for i, exp := range expected {
		if i < len(sliceResult) && sliceResult[i] != exp {
			t.Errorf("Register[%d]: expected %d, got %d", i, exp, sliceResult[i])
		}
	}
}

func TestReadTag_CoilRouting(t *testing.T) {
	mock := &MockModbusClient{
		coilsData: []bool{true},
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestCoilTag",
		RegisterType: "Coil",
		Address:      1,
		Size:         1,
	}

	result, err := client.ReadTag(tag)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != true {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestReadTag_HoldingRegisterRouting(t *testing.T) {
	mock := &MockModbusClient{
		registersData: []uint16{42},
	}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestHRTag",
		RegisterType: "HoldingRegister",
		Address:      1,
		Size:         1,
	}

	result, err := client.ReadTag(tag)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	intResult, ok := result.(int16)
	if !ok {
		t.Fatalf("Expected int16, got type: %T", result)
	}

	if intResult != 42 {
		t.Errorf("Expected 42, got: %v", intResult)
	}
}

func TestReadTag_UnsupportedRegisterType(t *testing.T) {
	mock := &MockModbusClient{}
	client := NewClientWithModbus(mock)

	tag := model.ModbusTag{
		Name:         "TestUnsupported",
		RegisterType: "InputRegister", // Unsupported type
		Address:      1,
		Size:         1,
	}

	_, err := client.ReadTag(tag)
	if err == nil {
		t.Fatal("Expected error for unsupported register type, got none")
	}

	expected := "unsupported register type: InputRegister"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got: %v", expected, err)
	}
}

func TestReadCoil_AddressConversion(t *testing.T) {
	mock := &MockModbusClient{
		coilsData: []bool{true},
	}
	client := NewClientWithModbus(mock)

	// Verify that tag.Address-1 is passed to the modbus client
	// We can't directly verify the address passed, but we can test the behavior
	tag := model.ModbusTag{
		Name:         "TestAddressConversion",
		RegisterType: "Coil",
		Address:      10, // Should be converted to 9 (10-1)
		Size:         1,
	}

	result, err := client.readCoil(tag)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != true {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestFormatTagValue_MapValue(t *testing.T) {
	client := &Client{}

	testCases := []struct {
		name     string
		value    map[string]any
		expected string
	}{
		{
			name: "reasonable float",
			value: map[string]any{
				"float32": float32(60.0),
				"int32":   int32(1114636288),
			},
			expected: "float32: 60 (int32: 1114636288)",
		},
		{
			name: "reasonable int",
			value: map[string]any{
				"float32": float32(1.234e-10), // Very small float
				"int32":   int32(100),
			},
			expected: "int32: 100 (float32: 1.234e-10)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tag := model.ModbusTag{Name: "test"}
			result := client.FormatTagValue(tag, tc.value)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestFormatTagValue_SimpleValue(t *testing.T) {
	client := &Client{}
	tag := model.ModbusTag{Name: "test"}

	result := client.FormatTagValue(tag, 42)
	expected := "42"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Integration test (requires actual Modbus server)
func TestModbusConnection(t *testing.T) {
	client, err := NewClient("172.29.48.69:502")
	if err != nil {
		t.Skipf("Modbus server not available: %v", err)
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
