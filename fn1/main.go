package webcam

import (
	"bytes"
	"fmt"
	_ "image/jpeg"
	"signature"

	pigo "github.com/esimov/pigo/core"

	_ "embed"
)

//go:embed facefinder
var cascade_facefinder []byte

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

	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,

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
	dets = classifier.ClusterDetections(dets, 0.2)

	ctx.Status = fmt.Sprintf("Detected %d", len(dets))

	return signature.Next(ctx)
}
