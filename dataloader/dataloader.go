package dataloader

import (
	"github.com/pkg/errors"
	"mtggokits/datacontainer"
	"mtggokits/dataloader/streamer"
)

func Get(name string, key datacontainer.MapKey) (interface{}, error) {
	s, ok := streamer.DataStreamers[name]
	if !ok {
		return nil, errors.New("not found streamer[" + name + "]")
	}
	c := s.GetContainer()
	if c == nil {
		return nil, errors.New("contain is nil, streamer[" + name + "]")
	}
	return c.Get(key)
}
