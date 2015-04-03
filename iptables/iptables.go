package iptables

import (
	"bufio"
	"bytes"
	//"fmt"
	"regexp"
	"log"
	"sort"
	"os/exec"
	"strconv"
)

type IPTable struct {
	Name	string
	Chains	[]IPTableChain
}

type IPTableChain struct {
	Name	string
	Rows	[]IPTableRow
}

type IPTableRow struct {
	Target		string // DNAT?
	Protocol	string
	SourceAddr	string
	Destination	string
	SpecDestIP	string
	SpecDestPort	int
	SpecSrcPort	int
}

func (table *IPTable) Load(tableName string) {
	table.Name = tableName
	out, err := exec.Command("iptables", "-t", tableName, "-L").Output()
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(out)
	scanner := bufio.NewScanner(reader)

	chainLine := regexp.MustCompile(`^Chain\s([^\s]+)`)
	ruleLine := regexp.MustCompile(`^(DNAT|MASQUERADE)\s+(all|tcp|udp|icmp)\s{2}--\s{2}([^\s]+)\s+([^\s]+)\s+(.*)$`)
	spec := regexp.MustCompile(`(all|tcp|udp|icmp)\sdpt:(\d{2,6})\sto\:([^\:]+)\:(\d{2,6})`)

	var chain *IPTableChain
	chain = nil
	for scanner.Scan() {
		line := scanner.Text()
		r := chainLine.FindStringSubmatch(line)

		if len(r) > 0 {
			chainName := r[len(r)-1]
			chain = table.AddChain(chainName)
			scanner.Scan()
		}

		s := ruleLine.FindStringSubmatch(line)
		if len(s) > 0 {
			target := s[1]
			protocol := s[2]
			source := s[3]
			destination := s[4]

			row := IPTableRow{}
			row.Target	= target
			row.Protocol	= protocol
			row.SourceAddr	= source
			row.Destination = destination

			if len(s) == 6 {
				sm := spec.FindStringSubmatch(s[5])
				if len(sm) > 0 {
					row.SpecDestIP = sm[3]
					row.SpecDestPort, _ = strconv.Atoi(sm[4])
					row.SpecSrcPort, _ = strconv.Atoi(sm[2])
				}
			}

			chain.Rows = append([]IPTableRow{row}, chain.Rows...)
		}

	}

}

func (table *IPTable) AddChain(chainName string) *IPTableChain {
	chain := table.FindChain(chainName)
	if chain != nil {
		return chain
	}
	newChain := IPTableChain{}
	newChain.Name = chainName

	table.Chains = append([]IPTableChain{newChain}, table.Chains...)
	return &newChain
}

func (table IPTable) FindChain(chainName string) *IPTableChain {
	i := sort.Search(len(table.Chains), func(i int) bool { return table.Chains[i].Name == chainName })
	if i < len(table.Chains) && table.Chains[i].Name == chainName {
		return &table.Chains[i]
	} else {
		return nil
	}
	return nil
}

func (table IPTable) AddRowToChain(chainName string, row IPTableRow) *IPTableChain {
	chain := table.FindChain(chainName)
	chain.Rows = append([]IPTableRow{row}, chain.Rows...)
	return chain
}

func (table IPTable) Commit() {
	//currentTable := GetNATTable()

}

func (chain IPTableChain) FindRow(protocol string, srcAddr string, destAddr string, srcPort int, destPort int) *IPTableRow {
	i := sort.Search(len(chain.Rows), func(i int) bool {
		res :=	chain.Rows[i].Protocol == protocol &&
			chain.Rows[i].SourceAddr == srcAddr &&
			chain.Rows[i].SpecDestIP == destAddr &&
			chain.Rows[i].SpecDestPort == destPort &&
			chain.Rows[i].SpecSrcPort == srcPort
		return res })

	if i < len(chain.Rows) {
		return &chain.Rows[i]
	} else {
		return nil
	}
}

func NewNATTable() IPTable {
	table := IPTable{}

	chain := IPTableChain{}
	chain.Name = "PREROUTING"

	table.Chains = append([]IPTableChain{chain}, table.Chains...)
	return table
}

func GetNATTable() IPTable {
	table := IPTable{}
	table.Load("nat")
	return table
}
