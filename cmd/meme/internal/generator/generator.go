package generator

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font/gofont/gobold"
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

	imgCtx := gg.NewContextForImage(inputImage)
	imgCtx.SetHexColor("#ff0") // YELLOW
	font, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return nil, errors.Wrap(err, "while parsing font")
	}
	fontFace := truetype.NewFace(font, &truetype.Options{Size: 16})
	imgCtx.SetFontFace(fontFace)
	imgCtx.DrawStringAnchored(quote, float64(inputImage.Bounds().Max.X/2.0), float64(inputImage.Bounds().Max.Y)*0.9, 0.5, 0.5)

	pReader, pWriter := io.Pipe()
	go func() {
		err := imgCtx.EncodePNG(pWriter)
		pWriter.CloseWithError(err)
	}()
	return pReader, nil
}

// dependencies
type QuoteProvider interface {
	Get() (string, error)
}
