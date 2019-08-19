// ******************************************************************************* //
//                                                                                 //
//                                                                                 //
// File: streamer_manager.go                                                       //
//                                                                                 //
// By: wangjia <jia.wang@mitegral.com>                                             //
//                                                                                 //
// Created: 2019/08/14 13:07:31 by wangjia                                         //
// Updated: 2019/08/14 13:07:31 by wangjia                                         //
//                                                                                 //
// ******************************************************************************* //

package streamer

import (
	"fmt"
	"github.com/pkg/errors"
)

var DataStreamers = make(map[string]DataStreamer)

type DataStreamerMangater struct {
}

func Register(name string, streamer DataStreamer) error {
	if _, ok := DataStreamers[name]; ok {
		fmt.Println("abcdefg")
		return errors.New("streamer[" + name + "] has already exist")
	}
	fmt.Println("Register: ", name)
	DataStreamers[name] = streamer
	return nil
}
