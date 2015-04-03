package iptables

import (
	"testing"
	"sort"
	"fmt"
)

func Log(t *testing.T, msg string) {
	t.Logf("\033[32m==> \033[0m%s", msg)
}

func TestLoad(t *testing.T) {
	Log(t, "Testing table load")
	table := &IPTable{}
	table.Load("nat")
	if table.Name != "nat" {
		t.Errorf("Table name not set correctly, it is %s", table.Name)
	}
	Log(t, fmt.Sprintf("Table name is %s", table.Name))

	i := sort.Search(len(table.Chains), func(i int) bool { return table.Chains[i].Name == "PREROUTING" })
	if i == 0 {
		t.Errorf("Table does not have a PREROUTING chain")
	}
}

func TestAddRow(t *testing.T) {
	table := &IPTable{}
	table.Load("nat")

	row := IPTableRow{}
	row.Target		= "DNAT"
	row.Protocol		= "all"
	row.SourceAddr	        = "10.10.10.0/24"
	row.Destination		= "anywhere"
	row.SpecDestIP		= "192.168.0.100"
	row.SpecDestPort	= 80
	row.SpecSrcPort		= 8080

	chain := table.AddRowToChain("PREROUTING", row)

	storedRow := chain.FindRow(row.Protocol,
		row.SourceAddr,
		row.SpecDestIP,
		row.SpecSrcPort,
		row.SpecDestPort)

	if storedRow == nil {
		t.Errorf("Couldn't retrieve new row")
	}

	Log(t, fmt.Sprintf("%+v\n", storedRow))
}

/*
row.SpecDestIP		= "192.168.0.100"

func TestCommit(t *testing.T) {
	table := &IPTable{}
	table.Load("nat")

	row := IPTableRow{}
	row.Target		= "DNAT"
	row.Protocol		= "all"
	row.SourceAddr	        = "192.168.0.0/24"
	row.Destination		= "anywhere"
	row.SpecDestIP		= "192.168.0.100"
	row.SpecDestPort	= 80
	row.SpecSrcPort		= 8080

	fmt.Printf("%+v", table)

	table.AddRowToChain("PREROUTING", row)

	fmt.Printf("%+v", table)
	// table.Commit()
}*/
