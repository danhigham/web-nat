package iptables

import (
	"fmt"
	"os"
)

type IPTable struct {
	Chains	[]IPTableChain
}

type IPTableChain struct {
	Rows	[]IPTableRow
}

type IPTableRow struct {
	DestIP 		string
	DestPort	int
	SrcPort		int
}

func (table IPTable) Load(tableName string) {
	
}


