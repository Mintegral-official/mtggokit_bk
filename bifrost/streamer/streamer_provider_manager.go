package streamer

import (
	"github.com/pkg/errors"
)

type StreamerProviderManager struct {
	StreamerProviders map[string]*StreamerProvider
}

func NewStreamerProviderManager() *StreamerProviderManager {
	return &StreamerProviderManager{
		StreamerProviders: make(map[string]*StreamerProvider),
	}
}

func (spm *StreamerProviderManager) RegiterProvider(name string, provider *StreamerProvider) error {
	if _, in := spm.StreamerProviders[name]; in {
		return errors.New("StreamerProvider[" + name + "] is already exist")
	}
	spm.StreamerProviders[name] = provider
	return nil
}

func (spm *StreamerProviderManager) GetProvider(name string, progress int64) *StreamerProvider {
	sp, in := spm.StreamerProviders[name]
	if in {
		return sp
	} else {
		return nil
	}
}
