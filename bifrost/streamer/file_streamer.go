package streamer

import (
	"bufio"
	"context"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	//"github.com/Mintegral-official/mtggokit/bifrost/log"
	"github.com/pkg/errors"
	"os"
	"time"
)

type LocalFileStreamer struct {
	container  container.Container
	cfg        *FileStreamerCfg
	scan       *bufio.Scanner
	updateMode UpdatMode
	//Logger     log.Logger
	hasInit bool
}

func NewFileStreamer(cfg *FileStreamerCfg) *LocalFileStreamer {
	fs := &LocalFileStreamer{
		cfg: cfg,
		//Logger: log,
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
	return fs.scan != nil && fs.scan.Scan()
}

func (fs *LocalFileStreamer) Next() (container.DataMode, container.MapKey, interface{}, error) {
	m, k, v, e := fs.cfg.DataParser.Parse([]byte(fs.scan.Text()), nil)
	return m, k, v, e
}

func (fs *LocalFileStreamer) UpdateData(ctx context.Context) error {
	if fs.cfg.IsSync {
		err := fs.updateData(ctx)
		if err != nil {
			return err
		}
	}
	go func() {
		for {
			inc := time.After(time.Duration(fs.cfg.Interval) * time.Second)
			select {
			case <-inc:
				if err := fs.updateData(ctx); err != nil {
					//log.New().Infof("streamer[%s] LoadInc error[%s]", fs.cfg.Name, err.Error())

				}
			}
		}
	}()
	return nil
}

func (fs *LocalFileStreamer) updateData(ctx context.Context) error {

	switch fs.updateMode {
	case Static, Dynamic:
		if fs.hasInit && fs.updateMode == Static {
			return nil
		}

		f, err := os.Open(fs.cfg.Path)
		defer func() { _ = f.Close() }()
		if err != nil {
			return err
		}
		fs.scan = bufio.NewScanner(f)
		return fs.container.LoadBase(fs)
	case Increment:
	case DynInc:
	default:
		return errors.New("not support mode[" + fs.updateMode.toString() + "]")
	}
	return nil
}
