package modbus

import (
	"encoding/binary"
	"fmt"
	"time"
	"unsafe"

	"opcmss/internal/model"

	"github.com/simonvetter/modbus"
)

type Client struct {
	client *modbus.ModbusClient
}

func NewClient(address string) (*Client, error) {
	c, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://" + address,
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	if err := c.Open(); err != nil {
		return nil, err
	}
	return &Client{client: c}, nil
}

func (c *Client) ReadTag(tag model.ModbusTag) (any, error) {
	switch tag.RegisterType {
	case "Coil":
		return c.readCoil(tag)
	case "HoldingRegister":
		return c.readHoldingRegister(tag)
	default:
		return nil, fmt.Errorf("unsupported register type: %s", tag.RegisterType)
	}
}

func (c *Client) readCoil(tag model.ModbusTag) (any, error) {
	// For coils, we read individual bits
	data, err := c.client.ReadCoils(tag.Address-1, tag.Size) // Modbus addresses are typically 1-based
	if err != nil {
		return nil, err
	}

	if tag.Size == 1 {
		return len(data) > 0 && data[0], nil
	}

	// For multiple coils, return slice of bools
	result := make([]bool, len(data))
	copy(result, data)
	return result, nil
}

func (c *Client) readHoldingRegister(tag model.ModbusTag) (any, error) {
	data, err := c.client.ReadRegisters(tag.Address-1, tag.Size, modbus.HOLDING_REGISTER)
	if err != nil {
		return nil, err
	}

	switch tag.Size {
	case 1:
		// Single 16-bit register - could be INT or BOOL
		if len(data) > 0 {
			return int16(data[0]), nil
		}
		return nil, fmt.Errorf("no data received")

	case 2:
		// Two 16-bit registers - could be REAL (float32) or DINT (int32)
		if len(data) < 2 {
			return nil, fmt.Errorf("insufficient data for 32-bit value")
		}

		// Convert to 32-bit value
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint16(bytes[0:2], data[0])
		binary.BigEndian.PutUint16(bytes[2:4], data[1])
		uint32Value := binary.BigEndian.Uint32(bytes)

		// Try to determine if it's float or int based on the value
		// For now, return both interpretations in a map
		return map[string]any{
			"int32":   int32(uint32Value),
			"float32": float32FromBits(uint32Value),
			"raw":     data,
		}, nil

	default:
		// Multiple registers, return as slice
		result := make([]int16, len(data))
		for i, val := range data {
			result[i] = int16(val)
		}
		return result, nil
	}
}

func (c *Client) Close() {
	c.client.Close()
}

func float32FromBits(bits uint32) float32 {
	return float32FromUint32(bits)
}

func float32FromUint32(bits uint32) float32 {
	return *(*float32)(unsafe.Pointer(&bits))
}
