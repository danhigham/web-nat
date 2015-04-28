package main


type TestData struct { }

func (t TestData) Containers() string {
	data := ` 
	[{
		"state": {
			"started_at": "2014-09-12T00:48:23.824260519Z",
			"pid": 5845,
			"running": true
		}, 
		"ports": [
		{
			"container_port": 8080,
			"port": 49159,
			"proto": "tcp"
		}
		],
		"engine": {
			"labels": [
			"local",
			"dev"
			],
			"memory": 4096,
			"cpus": 4,
			"addr": "http://172.16.1.50:2375",
			"id": "local"
		},
		"image": {
			"restart_policy": {},
			"labels": [
			""
			],
			"type": "service",
			"hostname": "cbe68bf32f1a",
			"environment": {
				"GOROOT": "/goroot",
				"GOPATH": "/gopath"
			},
			"memory": 256,
			"cpus": 0.08,
			"name": "ehazlett/go-demo:latest"
		},
		"id": "cbe68bf32f1a08218693dbee9c66ea018c1a99c75c463a76b"
	},
	{
		"state": {
			"started_at": "2014-09-12T00:48:23.824260519Z",
			"pid": 5846,
			"running": true
		}, 
		"ports": [
		{
			"container_port": 8080,
			"port": 49158,
			"proto": "tcp"
		}
		],
		"engine": {
			"labels": [
			"local",
			"dev"
			],
			"memory": 4096,
			"cpus": 4,
			"addr": "http://172.16.1.50:2375",
			"id": "local"
		},
		"image": {
			"restart_policy": {},
			"labels": [
			""
			],
			"type": "service",
			"hostname": "eca254ecd76e",
			"environment": {
				"GOROOT": "/goroot",
				"GOPATH": "/gopath"
			},
			"memory": 256,
			"cpus": 0.08,
			"name": "ehazlett/go-demo:latest"
		},
		"id": "eca254ecd76eb9d887995114ff811cc5b7c14fe13630"
	}]`

	return data
}

func (t TestData) IPTable() string {
	data := `[
		{ "id": "Wordpress 1234", "host": "host 1", "port_nat": [ { "container_port": 80, "port": 8080, "from": "192.168.3.0/24"} ] },
		{ "id": "Minecraft Server", "host": "host 2", "port_nat": [ { "container_port": 1234, "port": 1234, "from": "anywhere"} ] },
		{ "id": "MySQL", "host": "host 1", "port_nat": [ { "container_port": 3306, "port": 3306, "from": "192.168.3.0/24"} ] }
	]`

	return data
}

