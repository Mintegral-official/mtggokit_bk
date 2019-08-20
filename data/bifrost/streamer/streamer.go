package streamer

import (
	"github.com/Mintegral-official/mtggokit/data/container"
)

type DataStreamer interface {
	SetContainer(container.Container)
	GetContainer() container.Container

	UpdateData() error
}
