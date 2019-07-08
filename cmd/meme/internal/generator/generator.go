package generator

import (
	"fmt"
	"github.com/jpoz/gomeme"
	"github.com/pkg/errors"
	"image"
	"io"
	"math/rand"
	"os"
	"time"
)

type MemGen struct {
	quoteProvider QuoteProvider
	randomSource  *rand.Rand
}

func New(quoteProvider QuoteProvider) *MemGen {
	return &MemGen{quoteProvider: quoteProvider,
		randomSource: rand.New(rand.NewSource(int64(time.Now().Nanosecond())))}
}

func (g *MemGen) Get() (io.Reader, error) {
	quote, err := g.quoteProvider.Get()
	if err != nil {
		return nil, errors.Wrap(err, "while getting quote for meme")
	}

	imgPath := fmt.Sprintf("cmd/meme/resources/face-%d.jpg", g.randomSource.Intn(7))
	file, err := os.Open(imgPath)
	if err != nil {
		return nil, errors.Wrapf(err, "while opening image with path: [%s]", imgPath)
	}
	defer file.Close()

	inputImage, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.Wrapf(err, "while decoding image [%s]", imgPath)
	}

	config := gomeme.NewConfig()
	config.BottomText = quote
	meme := &gomeme.Meme{
		Config:   config,
		Memeable: gomeme.JPEG{Image: inputImage},
	}

	pReader, pWriter := io.Pipe()
	go func() {
		err := meme.Write(pWriter)
		pWriter.CloseWithError(err)
	}()
	return pReader, nil
}

// dependencies
type QuoteProvider interface {
	Get() (string, error)
}
