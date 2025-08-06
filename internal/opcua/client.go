package opcua

import (
	"context"
	"fmt"
	"time"

	"opcmss/internal/model"

	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
)

type Client struct {
	client *client.Client
	ctx    context.Context
}

func NewClient(endpoint string) (*Client, error) {
	ctx := context.Background()

	client, err := client.Dial(ctx, endpoint, client.WithInsecureSkipVerify())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OPC UA server: %w", err)
	}

	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

func (c *Client) ReadTag(tag model.OPCTag) (any, error) {
	node := ua.ParseNodeID(tag.NodeID)

	req := &ua.ReadRequest{
		NodesToRead: []ua.ReadValueID{
			{
				NodeID:      node,
				AttributeID: ua.AttributeIDValue,
			},
		},
	}
	val, err := c.client.Read(c.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	if len(val.Results) == 0 {
		return nil, fmt.Errorf("no results returned")
	}

	result := val.Results[0]
	if !result.StatusCode.IsGood() {
		return nil, fmt.Errorf("read failed with status: %v", result.StatusCode)
	}

	// Use result.Value directly (no .Value() method needed)
	actualValue := result.Value

	switch tag.DataType {
	case "BOOL":
		if b, ok := actualValue.(bool); ok {
			return b, nil
		}
		return false, fmt.Errorf("failed to convert value to bool, got type %T", actualValue)
	case "INT":
		// Handle different integer types that might be returned
		switch v := actualValue.(type) {
		case int16:
			return v, nil
		case int32:
			return int16(v), nil
		case int:
			return int16(v), nil
		case uint16:
			return int16(v), nil
		case uint32:
			return int16(v), nil
		default:
			return int16(0), fmt.Errorf("failed to convert value to int16, got type %T with value %v", actualValue, actualValue)
		}
	case "REAL":
		// Handle different numeric types that might be returned for REAL
		switch v := actualValue.(type) {
		case float32:
			return v, nil
		case float64:
			return float32(v), nil
		case int16:
			return float32(v), nil
		case int32:
			return float32(v), nil
		case int:
			return float32(v), nil
		case uint16:
			return float32(v), nil
		case uint32:
			return float32(v), nil
		default:
			return float32(0), fmt.Errorf("failed to convert value to float32, got type %T with value %v", actualValue, actualValue)
		}
	default:
		// Return the raw value for debugging
		return actualValue, nil
	}
}

func (c *Client) WriteTag(tag model.OPCTag, value any) error {
	node := ua.ParseNodeID(tag.NodeID)

	req := &ua.WriteRequest{
		NodesToWrite: []ua.WriteValue{
			{
				NodeID:      node,
				AttributeID: ua.AttributeIDValue,
				Value: ua.DataValue{
					Value: value,
				},
			},
		},
	}

	res, err := c.client.Write(c.ctx, req)
	if err != nil {
		return fmt.Errorf("error writing to node %s: %w", tag.NodeID, err)
	}
	if len(res.Results) > 0 && res.Results[0].IsGood() {
		return nil
	}
	return fmt.Errorf("failed to write to node %s: %s", tag.NodeID, res.Results[0])
}

func (c *Client) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c.client.Close(ctx)
}
