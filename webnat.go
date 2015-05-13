package main

import (
	"net/url"
	"fmt"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/danhigham/webnat/iptables"
)

type ShipyardContainer struct {
	ID		string						`json:"id"`
	State	ShipyardContainerState		`json:"state"`
	Ports	[]ShipyardContainerPort		`json:"ports"`
	Engine	ShipyardContainerEngine		`json:"engine"`
	Image	ShipyardContainerImage		`json:"image"`
	PortNat []TranslatedPort			`json:"port_nat"`
}

type ShipyardContainerState struct {
	StartedAt	string					`json:"started_at"`
	Pid			int						`json:"pid"`
	Running		bool					`json:"running"`
}

type ShipyardContainerPort struct {
	ContainerPort	int					`json:"container_port"`
	Port			int					`json:"port"`
	Protocol		string				`json:"proto"`
}

type ShipyardContainerEngine struct {
	Memory			int					`json:"memory"`
	CPUs			int					`json:"cpus"`
	Address			string				`json:"addr"`
	ID				string				`json:"id"`
}

type ShipyardContainerImage struct {
	Type			string				`json:"type"`
	Hostname		string				`json:"hostname"`
	Memory			int					`json:"memory"`
	CPUs			float64				`json:"cpus"`
	Name			string				`json:"name"`
}

type TranslatedPort struct {
	ContainerPort	int					`json:"container_port"`
	Port			int					`json:"port"`
	From			string				`json:"from"`
}

func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")

	r.GET("/shipyard_containers", func(c *gin.Context) {
		res := TestData{}.Containers()
		writeStringAsJSON(c, 200, res)
	})

	r.GET("/ip_table", func(c *gin.Context) {
		res := TestData{}.IPTable()
		writeStringAsJSON(c, 200, res)
	})

	r.GET("/ip_table_raw", func(c *gin.Context) {
		res := iptables.GetNATTable()
		c.JSON(200, res)
	})

	r.GET("/graph_data", func(c *gin.Context) {
		var containers []ShipyardContainer
		//table := iptables.GetNATTable()
		ret := TestData{}.Containers()
		err := json.Unmarshal([]byte(ret), &containers)
		if err != nil {
			fmt.Println("error:", err)
		}
		for c := range containers {
			container := containers[c]
			uri, _ := url.Parse(container.Engine.Address)
			fmt.Println(uri.Host)
		}

		c.JSON(200, containers)
	})

	r.Run(":8080")
}

func writeStringAsJSON(c *gin.Context, status int, body string) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte(body))
	c.Writer.Flush()
}
