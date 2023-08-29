package webcam

import (
	"bytes"
	"fmt"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"math"
	"signature"

	pigo "github.com/esimov/pigo/core"
	"github.com/fogleman/gg"

	_ "embed"
)

//go:embed facefinder
var cascade_facefinder []byte

var dc *gg.Context

func Scale(ctx *signature.Context) (*signature.Context, error) {
	ctx.Status = "Hello"

	data := bytes.NewReader([]byte(ctx.Frame))
	src, err := pigo.DecodeImage(data)
	if err != nil {
		ctx.Status = fmt.Sprintf("Cannot open the image file: %v", err)
		return ctx, err
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	// Put the image in the gg context so we can manipulate it...
	dc = gg.NewContext(cols, rows)
	dc.DrawImage(src, 0, 0)

	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     1000,
		ShiftFactor: 0.15,
		ScaleFactor: 1.15,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	pigo := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := pigo.Unpack(cascade_facefinder)
	if err != nil {
		ctx.Status = fmt.Sprintf("Error reading the cascade file: %s", err)
		return ctx, err
	}

	angle := 0.0 // cascade rotation angle. 0.0 is 0 radians and 1.0 is 2*pi radians

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0.15)

	o := "Detection finished"

	qThresh := float32(5.0)

	detections := make([]signature.Detection, 0)

	for _, d := range dets {
		if d.Q > qThresh {
			o = fmt.Sprintf("%s [%d,%d,%d, %.2f]", o, d.Col, d.Row, d.Scale, d.Q)

			dc.DrawArc(
				float64(d.Col),
				float64(d.Row),
				float64(d.Scale/2),
				0,
				2*math.Pi,
			)

			dc.SetLineWidth(2.0)
			dc.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
			dc.Stroke()

			detections = append(detections, signature.Detection{
				Row:   uint32(d.Row),
				Col:   uint32(d.Col),
				Scale: uint32(d.Scale),
			})
		}
	}

	ctx.Detections = detections
	ctx.Status = o

	// Encode the image with faces detected...
	img := dc.Image()
	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, img, &jpeg.Options{Quality: 100})

	ctx.Frame = buff.Bytes()

	return signature.Next(ctx)
}
