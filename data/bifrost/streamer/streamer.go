package streamer

import (
	"mtggokits/data/container"
)

type DataStreamer interface {
	SetContainer(container.Container)
	GetContainer() container.Container

	UpdateData() error
}
