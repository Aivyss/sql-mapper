package component

import "encoding/xml"

type AppCtxComponent struct {
	xml.Name              `xml:"Context"`
	QueryClientsComponent `xml:"QueryClients"`
}

type QueryClientsComponent struct {
	clients []QueryClientComponent `xml:"QueryClient"`
}

type QueryClientComponent struct {
	identifier string `xml:"identifier,attr"`
	filePath   string `xml:"filePath,attr"`
}
