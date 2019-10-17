package streamer

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/pkg/errors"
	"os"
	"time"
)

type LocalFileStreamer struct {
	container container.Container
	cfg       *LocalFileStreamerCfg
	scan      *bufio.Scanner
	result    []ParserResult
	curLen    int
	hasInit   bool
	modTime   time.Time
	totalNum  int64
	errorNum  int64
	startTime int64
	endTime   int64
}

func NewFileStreamer(cfg *LocalFileStreamerCfg) *LocalFileStreamer {
	fs := &LocalFileStreamer{
		cfg: cfg,
	}
	return fs
}

func (fs *LocalFileStreamer) SetContainer(container container.Container) {
	fs.container = container
}

func (fs *LocalFileStreamer) GetContainer() container.Container {
	return fs.container
}

func (fs *LocalFileStreamer) GetSchedInfo() *SchedInfo {
	return &SchedInfo{
		TimeInterval: fs.cfg.Interval,
	}
}

func (fs *LocalFileStreamer) HasNext() bool {
	return fs.curLen < len(fs.result) || fs.scan != nil && fs.scan.Scan()
}

func (fs *LocalFileStreamer) Next() (container.DataMode, container.MapKey, interface{}, error) {
	fs.totalNum++
	if fs.curLen < len(fs.result) {
		r := fs.result[fs.curLen]
		fs.curLen++
		if r.Err != nil {
			fs.errorNum++
		}
		return r.DataMode, r.Key, r.Value, r.Err
	}
	result := fs.cfg.DataParser.Parse([]byte(fs.scan.Text()), nil)
	if result == nil {
		fs.errorNum++
		return container.DataModeAdd, nil, nil, errors.New(fmt.Sprintf("Parser error"))
	}
	fs.curLen = 0
	fs.result = result
	if fs.curLen < len(fs.result) {
		r := fs.result[fs.curLen]
		fs.curLen++
		if r.Err != nil {
			fs.errorNum++
		}
		return r.DataMode, r.Key, r.Value, r.Err
	}
	fs.errorNum++
	return container.DataModeAdd, nil, nil, errors.New(fmt.Sprintf("Index[%d] error, len[%d]", fs.curLen, len(fs.result)))
}

func (fs *LocalFileStreamer) UpdateData(ctx context.Context) error {
	if fs.cfg.IsSync {
		fs.startTime = time.Now().UnixNano()
		err := fs.updateData(ctx)
		fs.endTime = time.Now().UnixNano()
		if err != nil {
			fs.WarnStatus("LoadBase error: " + err.Error())
			return err
		}
		fs.InfoStatus("LoadBase succ")
	}
	go func() {
		for {
			inc := time.After(time.Duration(fs.cfg.Interval) * time.Second)
			select {
			case <-ctx.Done():
				return
			case <-inc:
				fs.startTime = time.Now().UnixNano()
				err := fs.updateData(ctx)
				fs.endTime = time.Now().UnixNano()
				if err != nil {
					fs.WarnStatus("LoadBase error: " + err.Error())
				} else {
					fs.InfoStatus("LoadBase succ")
				}
			}
		}
	}()
	return nil
}

func (fs *LocalFileStreamer) updateData(ctx context.Context) error {

	switch fs.cfg.UpdatMode {
	case Static, Dynamic:
		fs.totalNum = 0
		fs.errorNum = 0
		if fs.hasInit && fs.cfg.UpdatMode == Static {
			return nil
		}

		f, err := os.Open(fs.cfg.Path)
		defer func() { _ = f.Close() }()
		if err != nil {
			return err
		}
		stat, _ := f.Stat()
		modTime := stat.ModTime()
		if modTime.After(fs.modTime) {
			fs.modTime = modTime
			fs.scan = bufio.NewScanner(f)
			return fs.container.LoadBase(fs)
		}
	case Increment:
	case DynInc:
	default:
		return errors.New("not support mode[" + fs.cfg.UpdatMode.toString() + "]")
	}
	return nil
}

func (ms *LocalFileStreamer) InfoStatus(s string) {
	if ms.cfg.Logger != nil {
		ms.cfg.Logger.Infof("streamer[%s] %s, totalNum[%d], errorNum[%d], timeUsed[%d]", ms.cfg.Name, s, ms.totalNum, ms.errorNum, (ms.endTime-ms.startTime)/10e6)
	}
}

func (ms *LocalFileStreamer) WarnStatus(s string) {
	if ms.cfg.Logger != nil {
		ms.cfg.Logger.Warnf("streamer[%s] %s, totalNum[%d], errorNum[%d], timeUsed[%d]", ms.cfg.Name, s, ms.totalNum, ms.errorNum, (ms.endTime-ms.startTime)/10e6)
	}
}
