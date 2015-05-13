package iptables

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"log"
	"os/exec"
	"strconv"
	"github.com/olekukonko/tablewriter"
)

type IPTable struct {
	Name	string			`json:"name"`
	Chains	[]IPTableChain  `json:"chains"`
}

type IPTableChain struct {
	Name	string			`json:"name"`
	Rows	[]IPTableRow	`json:"rows"`
}

type IPTableRow struct {
	Index		int			`json:"index"`
	Target		string		`json:"target"`
	Protocol	string		`json:"protocol"`
	SourceAddr	string		`json:"source_address"`
	Destination	string		`json:"destination"`
	SpecDestIP	string		`json:"spec_destination_ip"`
	SpecDestPort	int		`json:"spec_destination_port"`
	SpecSrcPort	int			`json:"spec_source_port"`
}

func (row IPTableRow) ToArray() []string {
	return []string {
		row.Target,
		row.Protocol,
		row.SourceAddr,
		row.Destination,
		row.SpecDestIP,
		strconv.Itoa(row.SpecDestPort),
		strconv.Itoa(row.SpecSrcPort) }
}

func (table *IPTable) Load(tableName string) {
	table.Name = tableName
	table.Chains = []IPTableChain{}

	out, err := exec.Command("iptables", "-t", tableName, "-L", "-n").Output()
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

			if chain != nil {
				table.AddRowToChain(chain.Name, row)
			}
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

func (table *IPTable) FindChain(chainName string) *IPTableChain {
	x := -1

	for i := range table.Chains {
		if table.Chains[i].Name == chainName {
			x = i
		}
	}

	if x > -1 {
		return &table.Chains[x]
	} else {
		return nil
	}
}

func (table IPTable) AddRowToChain(chainName string, row IPTableRow) *IPTableChain {
	chain := table.FindChain(chainName)
	chain.Rows = append(chain.Rows, row)
	return chain
}

func (table IPTable) Commit() {
	currentTable := IPTable{}
	currentTable.Load(table.Name)

	//iterate throrugh chains and table rows
	for c := range table.Chains {
		chain := table.Chains[c]
		for r := range chain.Rows {
			row := chain.Rows[r]

			currentChain := currentTable.FindChain(chain.Name)
			currentRow := currentChain.FindRow(row.Protocol, row.SourceAddr, row.SpecDestIP,
				row.SpecSrcPort, row.SpecDestPort)

			if currentRow == nil {
				// add rule to table
				row.Commit(table.Name, chain.Name)
			}
		}
	}

	for c := range currentTable.Chains {
		chain := currentTable.Chains [c]

		for r := range chain.Rows {
			row := chain.Rows[r]

			newChain := table.FindChain(chain.Name)
			newRow := newChain.FindRow(row.Protocol, row.SourceAddr, row.SpecDestIP,
				row.SpecSrcPort, row.SpecDestPort)

			if newRow == nil {
				cmd := exec.Command("iptables", "-t", table.Name, "-D", chain.Name, strconv.Itoa(r + 1))
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}

func (table IPTable) Dump() string {
	b := new (bytes.Buffer)
	w := bufio.NewWriter(b)

	for _, c := range table.Chains {
		if len(c.Rows) > 0 {
			fmt.Fprintf(w, "\n---=== %s ===---\n%s", c.Name, c.ToTable())
		}
	}

	w.Flush()
	return fmt.Sprint(b)
}

func (chain IPTableChain) FindRow(protocol string, srcAddr string, destAddr string, srcPort int, destPort int) *IPTableRow {
	x := -1
	for i := range chain.Rows {
		chain := chain.Rows[i]
		if chain.Protocol == protocol &&
		   chain.SourceAddr == srcAddr &&
		   chain.SpecDestIP == destAddr &&
		   chain.SpecDestPort == destPort &&
		   chain.SpecSrcPort == srcPort {
			x = i
		}
	}

	if x > -1  {
		row := &chain.Rows[x]
		row.Index = x
		return row
	} else {
		return nil
	}
}

func (chain IPTableChain) ToTable() string {
	b := new(bytes.Buffer)
	w := bufio.NewWriter(b)

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Target", "Protocol", "Source Addr", "Destination", "Dest IP", "Dest Port", "Src Port"})

	for _, r := range chain.Rows {
		table.Append(r.ToArray())
	}

	table.Render()
	w.Flush()
	return fmt.Sprintf("\n%s", b)
}

func (row IPTableRow) Commit(tableName string, chainName string) {
	var cmd *exec.Cmd

	if row.SpecDestIP == "" && row.SpecDestPort == 0 && row.SpecSrcPort == 0 {
		cmd = exec.Command("iptables", "-t", tableName, "-A", chainName, "-j", row.Target, "-p", row.Protocol)
		// iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
	} else if row.SourceAddr == "" {
		// iptables -t nat -A PREROUTING -p tcp --dport 1111 -j DNAT --to-destination 2.2.2.2:1111
		cmd = exec.Command("iptables", "-t", tableName, "-A", chainName,
			"-p", row.Protocol, "--dport", strconv.Itoa(row.SpecSrcPort), "-j", "DNAT",
			"--to-destination", fmt.Sprintf("%s:%s", row.SpecDestIP, strconv.Itoa(row.SpecDestPort)))
	} else {
		// iptables -t nat -A PREROUTING -s 192.168.1.1 -p tcp --dport 1111 -j DNAT --to-destination 2.2.2.2:1111
		cmd = exec.Command("iptables", "-t", tableName, "-A", chainName, "-s", row.SourceAddr,
			"-p", row.Protocol, "--dport", strconv.Itoa(row.SpecSrcPort), "-j", "DNAT",
			"--to-destination", fmt.Sprintf("%s:%s", row.SpecDestIP, strconv.Itoa(row.SpecDestPort)))
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func (chain *IPTableChain) RemoveRow(i int) {
	chain.Rows = append(chain.Rows[:i], chain.Rows[i+1:]...)
}

func GetNATTable() IPTable {
	table := IPTable{}
	table.Load("nat")
	return table
}
