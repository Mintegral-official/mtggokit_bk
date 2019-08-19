// ******************************************************************************* //
//                                                                                 //
//                                                                                 //
// File: streamer.go                                                               //
//                                                                                 //
// By: wangjia <jia.wang@mitegral.com>                                             //
//                                                                                 //
// Created: 2019/08/14 12:28:02 by wangjia                                         //
// Updated: 2019/08/14 12:28:02 by wangjia                                         //
//                                                                                 //
// ******************************************************************************* //

package streamer

import (
	"mtggokits/data/container"
)

type DataStreamer interface {
	SetContainer(container.Container)
	GetContainer() container.Container

	UpdateData() error
}
