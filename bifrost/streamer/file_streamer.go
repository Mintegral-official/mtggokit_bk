package streamer

import (
	"bufio"
	"context"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/pkg/errors"
	"os"
)

type FileStreamer struct {
	container  container.Container
	cfg        *FileStreamerCfg
	f          *os.File
	scan       *bufio.Scanner
	updateMode UpdatMode
}

func NewFileStreamer(cfg *FileStreamerCfg) (*FileStreamer, error) {
	fs := &FileStreamer{
		cfg: cfg,
	}
	return fs, nil
}

func (fs *FileStreamer) SetContainer(container container.Container) {
	fs.container = container
}

func (fs *FileStreamer) GetContainer() container.Container {
	return fs.container
}

func (fs *FileStreamer) HasNext() bool {
	return fs.scan.Scan()
}

func (fs *FileStreamer) Next() (container.DataMode, container.MapKey, interface{}, error) {
	m, k, v, e := fs.cfg.DataParser.Parse([]byte(fs.scan.Text()), nil)
	return m, k, v, e
}

func (fs *FileStreamer) UpdateData(ctx context.Context) error {
	switch fs.updateMode {
	case Static:
	case Dynamic:
		if fs.f != nil {
			_ = fs.f.Close()
		}
		f, err := os.Open(fs.cfg.Path)
		if err != nil {
			return err
		}
		fs.f = f
		_, _ = fs.f.Seek(0, 0)
		return fs.container.LoadBase(fs)
	case Increase:
	case DynInc:
	default:
		return errors.New("not support mode[" + fs.updateMode.toString() + "]")
	}
	return nil

}
