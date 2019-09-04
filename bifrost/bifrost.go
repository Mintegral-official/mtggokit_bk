package bifrost

import (
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/Mintegral-official/mtggokit/bifrost/streamer"
	"github.com/pkg/errors"
)

// Logger for log
type Logger interface {
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

type Bifrost struct {
	DataStreamers map[string]streamer.DataStreamer
	logger        *Logger
}

func NewBifrost() *Bifrost {
	return &Bifrost{
		DataStreamers: make(map[string]streamer.DataStreamer),
	}
}

func (l *Bifrost) Get(name string, key container.MapKey) (interface{}, error) {
	s, ok := l.DataStreamers[name]
	if !ok {
		return nil, errors.New("not found streamer[" + name + "]")
	}
	c := s.GetContainer()
	if c == nil {
		return nil, errors.New("contain is nil, streamer[" + name + "]")
	}
	return c.Get(key)
}

func (l *Bifrost) Register(name string, streamer streamer.DataStreamer) error {
	if _, ok := l.DataStreamers[name]; ok {
		return errors.New("streamer[" + name + "] has already exist")
	}
	l.DataStreamers[name] = streamer
	return nil
}

func (l *Bifrost) GetStreamer(name string) (streamer.DataStreamer, error) {
	s, ok := l.DataStreamers[name]
	if !ok {
		return nil, errors.New("not found streamer[" + name + "]")
	}
	return s, nil
}
