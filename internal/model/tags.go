package model

type ModbusTag struct {
	Name          string `json:"name"`
	RegisterType  string `json:"register_type"`  // "Coil" or "HoldingRegister"
	Address       uint16 `json:"address"`        // The modbus address
	ModbusAddress uint32 `json:"modbus_address"` // The full modbus address (e.g., 400002)
	Size          uint16 `json:"size"`           // Number of registers/coils
	Range         string `json:"range"`          // Range like "2..2" or "2210..2211"
}

type OPCTag struct {
	Name        string
	NodeID      string
	DataType    string
	Description string
}
