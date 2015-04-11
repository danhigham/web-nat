package iptables

import (
	"testing"
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

	Log(t, table.Dump())

}

func TestAddRow(t *testing.T) {
	table := &IPTable{}
	table.Load("nat")

	row := IPTableRow{}
	row.Target		= "DNAT"
	row.Protocol		= "tcp"
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

func TestCommit(t *testing.T) {
	table := &IPTable{}
	table.Load("nat")

	row := IPTableRow{}
	row.Target		= "DNAT"
	row.Protocol		= "tcp"
	row.SourceAddr	        = "192.168.0.0/24"
	row.Destination		= "anywhere"
	row.SpecDestIP		= "192.168.0.100"
	row.SpecDestPort	= 80
	row.SpecSrcPort		= 8080

	table.AddRowToChain("PREROUTING", row)

	table.Commit()

	newTable := IPTable{}
	newTable.Load("nat")

	chain := newTable.FindChain("PREROUTING")
	Log(t, chain.ToTable())
	newRow := chain.FindRow(row.Protocol, row.SourceAddr, row.SpecDestIP,
		row.SpecSrcPort, row.SpecDestPort)

	if newRow == nil {
		t.Errorf("Couldn't find added row")
	}

	chain.RemoveRow(newRow.Index)
	newTable.Commit()

	oldRow := chain.FindRow(row.Protocol, row.SourceAddr, row.SpecDestIP,
		row.SpecSrcPort, row.SpecDestPort)

	if oldRow != nil {
		t.Errorf("Old row still exists in table")
	}
}

func TestAddMasqueradeRule(t *testing.T){
	table := &IPTable{}
	table.Load("nat")

	chain := table.FindChain("POSTROUTING")
	if len(chain.Rows) > 0 {
		chain.RemoveRow(0)
		table.Commit()
	}

	if len(chain.Rows) > 0 {
		t.Errorf("Couldn't remove row")
	}

	chainLen := len(chain.Rows)

	row := IPTableRow{}
	row.Target = "MASQUERADE"
	row.Protocol = "all"
	row.SourceAddr = "anywhere"
	row.Destination = "anywhere"

	table.AddRowToChain("POSTROUTING", row)
	table.Commit()

	if len(chain.Rows) <= chainLen {
		t.Errorf("Couldn't add masqerade row back")
	}

	Log(t, chain.ToTable())
}
