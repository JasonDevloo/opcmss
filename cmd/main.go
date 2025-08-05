package main

import (
	"fmt"
	"log"

	"opcmss/internal/modbus"
	"opcmss/internal/parser"
)

func main() {
	tags, err := parser.ParseTagsTSV("/home/maimus/GoProjects/OPCvsModSymSrv/cmd/example_tags.tsv")
	if err != nil {
		log.Fatal(err)
	}

	client, err := modbus.NewClient("172.29.48.69:502")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	for _, tag := range tags {
		val, err := client.ReadTag(tag)
		if err != nil {
			log.Printf("Error reading %s: %v", tag.Name, err)
			continue
		}
		fmt.Printf("%s [%s] = %s\n", tag.Name, tag.RegisterType, client.FormatTagValue(tag, val))
	}
}
