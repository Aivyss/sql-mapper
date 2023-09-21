package component

import (
	"encoding/xml"
)

type AppCtxComponent struct {
	Name                 xml.Name              `xml:"Context"`
	QueryClientComponent QueryClientsComponent `xml:"QueryClients"`
}

type QueryClientsComponent struct {
	Clients []QueryClientComponent `xml:"QueryClient"`
}

type QueryClientComponent struct {
	Identifier string `xml:"identifier,attr"`
	FilePath   string `xml:"filePath,attr"`
}
