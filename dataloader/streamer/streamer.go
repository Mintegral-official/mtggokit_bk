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
	"context"
	"mtggokits/datacontainer"
)

type DataStreamer interface {
	SetContainer(datacontainer.Container)
	GetContainer() datacontainer.Container

	LoadBase(ctx context.Context) error
	LoadInc(ctx context.Context) error
}
