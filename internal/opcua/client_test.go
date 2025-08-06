package opcua

import (
	"testing"

	"opcmss/internal/model"
)

func TestOPCClient_ReadWrite(t *testing.T) {
	client, err := NewClient("opc.tcp://localhost:4840")
	if err != nil {
		t.Skipf("OPC server not available: %v", err)
	}
	defer client.Close()

	tag := model.OPCTag{
		Name:     "TestReal",
		NodeID:   "ns=2;s=Test.RealValue",
		DataType: "REAL",
	}

	writeVal := float32(23.5)
	err = client.WriteTag(tag, writeVal)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	readVal, err := client.ReadTag(tag)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	valFloat, ok := readVal.(float32)
	if !ok {
		t.Fatalf("Expected float32, got %T", readVal)
	}

	if valFloat != writeVal {
		t.Errorf("Mismatch: wrote %.2f, read %.2f", writeVal, valFloat)
	}
}
