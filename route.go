package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Route struct {
	Label        string
	TargetDomain string `json:"target-domain"`
	TargetPort   int    `json:"target-port"`
	ListenPort   int    `json:"listen-port"`
}

type LinkState struct {
	DelayTargetResponse  int  `json:"delay-target-response"`
	SendingToClient      bool `json:"sending-to-client"`
	ReceivingFromClient  bool `json:"receiving-from-client"`
	SendingToService     bool `json:"sending-to-service"`
	ReceivingFromService bool `json:"receiving-from-service"`
}

func readRoutes() (routes []Route, err error) {
	f, err := os.Open("routes.json")
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&routes)
	return
}

func (r Route) tcp() string {
	return fmt.Sprintf("%s:%d", r.TargetDomain, r.TargetPort)
}
