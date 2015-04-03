package iptables

import (
	"testing"
	"fmt"
)

func TestLoad(t *testing.T) {
	table := &IPTable{}
	table.Load("nat")
	fmt.Println(table.Name)

