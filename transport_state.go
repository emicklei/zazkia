/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"sync"
)

// TransportState specifies how the transport of data is affected.
type TransportState struct {
	mutex                        *sync.RWMutex // for atomically accessing values
	Verbose                      bool          `json:"verbose"`
	ThrottleServiceResponse      int           `json:"throttle-service-response"` // bytes per second
	DelayServiceResponse         int           `json:"delay-service-response"`    // milliseconds
	SendingToClient              bool          `json:"sending-to-client"`
	ReceivingFromClient          bool          `json:"receiving-from-client"`
	SendingToService             bool          `json:"sending-to-service"`
	ReceivingFromService         bool          `json:"receiving-from-service"`
	ServiceResponseCorruptMethod string        `json:"service-response-corrupt-method"`
	AcceptConnections            bool          `json:"accept-connections"`
}

func (t *TransportState) isAcceptConnections() bool {
	t.mutex.RLock()
	v := t.AcceptConnections
	t.mutex.RUnlock()
	return v
}

func (t *TransportState) toggleAcceptConnections() bool {
	t.mutex.Lock()
	v := !t.AcceptConnections
	t.AcceptConnections = v
	t.mutex.Unlock()
	return v
}

// TODO protect access to other values

func (s *TransportState) String() string {
	return fmt.Sprintf("accept=%v,delay=%d,s2c=%v,rfc=%v,s2s=%v,rfs=%v,thr=%d,cor=%s,v=%v",
		s.isAcceptConnections(),
		s.DelayServiceResponse, s.SendingToClient, s.ReceivingFromClient, s.SendingToService, s.ReceivingFromService,
		s.ThrottleServiceResponse, s.ServiceResponseCorruptMethod, s.Verbose)
}

func newDefaultTransportState() *TransportState {
	return &TransportState{
		mutex:                new(sync.RWMutex),
		SendingToClient:      true,
		ReceivingFromClient:  true,
		SendingToService:     true,
		ReceivingFromService: true,
		AcceptConnections:    true,
	}
}
