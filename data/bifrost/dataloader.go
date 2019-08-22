package bifrost

import (
	"github.com/Mintegral-official/mtggokit/data/bifrost/streamer"
	"github.com/Mintegral-official/mtggokit/data/container"
	"github.com/pkg/errors"
)

// Logger for log
type Logger interface {
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

type Loader struct {
	DataStreamers map[string]streamer.DataStreamer
	sched         Sched
	logger        *Logger
}

func NewLoader() *Loader {
	return &Loader{
		DataStreamers: make(map[string]streamer.DataStreamer),
	}
}

func (l *Loader) Get(name string, key container.MapKey) (interface{}, error) {
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

func (l *Loader) Register(name string, streamer streamer.DataStreamer) error {
	if _, ok := l.DataStreamers[name]; ok {
		return errors.New("streamer[" + name + "] has already exist")
	}
	l.DataStreamers[name] = streamer
	l.sched.AddStreamer(name, streamer)
	return nil
}

func (l *Loader) GetStreamer(name string) (streamer.DataStreamer, error) {
	s, ok := l.DataStreamers[name]
	if !ok {
		return nil, errors.New("not found streamer[" + name + "]")
	}
	return s, nil
}
