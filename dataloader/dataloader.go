package dataloader

import (
	"fmt"
	"github.com/pkg/errors"
	"mtggokits/datacontainer"
	"mtggokits/dataloader/streamer"
)

var DataStreamers = make(map[string]streamer.DataStreamer)

func Get(name string, key datacontainer.MapKey) (interface{}, error) {
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

func Register(name string, streamer streamer.DataStreamer) error {
	if _, ok := DataStreamers[name]; ok {
		fmt.Println("abcdefg")
		return errors.New("streamer[" + name + "] has already exist")
	}
	fmt.Println("Register: ", name)
	DataStreamers[name] = streamer
	return nil
}
