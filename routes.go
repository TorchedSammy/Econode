package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Route interface{
	Name() string
	SessionRequired() bool
	DataTransformer([]byte) (interface{}, error)
	Execute(client *Client, data interface{})
}

type SingleRoute struct{
	MethodName string
	Session bool
	//TransformHandler
	payloadType string
	rhandler reflect.Value
}

func createRoute(name, payload string, session bool, dataTransformer, handler interface{}) Route {
	return &SingleRoute{
		MethodName: name,
		Session: session,
		rhandler: reflect.ValueOf(handler),
		payloadType: payload,
	}
}

func (s *SingleRoute) Name() string {
	return s.MethodName
}

func (s *SingleRoute) SessionRequired() bool {
	return s.Session
}

func (s *SingleRoute) DataTransformer(data []byte) (transformed interface{}, err error) {
	if s.payloadType != "" {
		switch s.payloadType {
			case "auth":
				payload := AuthPayload{}
				err = json.Unmarshal(data, &payload)
				transformed = payload
		}
	}

	return
}

func (s *SingleRoute) Execute(c *Client, data interface{}) {
	s.rhandler.Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(data)})
}

func (e *Econetwork) getRoute(method string) Route {
	return e.routes[method]
}

func (e *Econetwork) addRoutes(routes []Route) {
	for _, rt := range routes {
		e.routes[rt.Name()] = rt
	}
}
