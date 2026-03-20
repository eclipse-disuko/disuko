// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package observermngmt

type EventCb func(EventId, interface{})

type ObserverManagement struct {
	cbs map[EventId][]EventCb
}

var defaultMngmt *ObserverManagement

func init() {
	defaultMngmt = &ObserverManagement{
		cbs: make(map[EventId][]EventCb),
	}
}

func Init() *ObserverManagement {
	return &ObserverManagement{
		cbs: make(map[EventId][]EventCb),
	}
}

func (o *ObserverManagement) RegisterHandler(id EventId, cb EventCb) {
	o.cbs[id] = append(o.cbs[id], cb)
}

func (o *ObserverManagement) FireEvent(id EventId, arg interface{}) {
	for _, cb := range o.cbs[id] {
		cb(id, arg)
	}
}

func RegisterHandler(id EventId, cb EventCb) {
	defaultMngmt.RegisterHandler(id, cb)
}

func FireEvent(id EventId, arg interface{}) {
	defaultMngmt.FireEvent(id, arg)
}
