package dataloader

import (
	"fmt"
	"github.com/pkg/errors"
	"mtggokits/data/container"
	streamer2 "mtggokits/data/dataloader/streamer"
)


type Loader struct {
	DataStreamers map[string]streamer2.DataStreamer
}

func NewLoader() *Loader {
	return Loader{DataStreamers:new(map[string]streamer2.DataStreamer))}
}

func(* loader) Get(name string, key container.MapKey) (interface{}, error) {
	s, ok := DataStreamers[name]
	if !ok {
		return nil, errors.New("not found streamer[" + name + "]")
	}
	c := s.GetContainer()
	if c == nil {
		return nil, errors.New("contain is nil, streamer[" + name + "]")
	}
	return c.Get(key)
}

func(* loader) Register(name string, streamer streamer2.DataStreamer) error {
	if _, ok := DataStreamers[name]; ok {
		return errors.New("streamer[" + name + "] has already exist")
	}
	DataStreamers[name] = streamer
	return nil
}
