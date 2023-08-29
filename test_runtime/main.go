package main

import (
	"context"
	"fmt"
	_ "image/jpeg"
	"os"
	sig "signature"

	scale "github.com/loopholelabs/scale"
	"github.com/loopholelabs/scale/scalefunc"
)

func main() {

	fmt.Printf("Starting... Using image %s\n", os.Args[1])

	fn, err := scalefunc.Read("./local-webcam-latest.scale")
	if err != nil {
		panic(err)
	}

	r, err := scale.New(context.Background(), sig.New, []*scalefunc.Schema{fn})
	if err != nil {
		panic(err)
	}

	i, err := r.Instance(nil)
	if err != nil {
		panic(err)
	}

	s := sig.New()
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	s.Context.Frame = data
	s.Context.Status = "Input"

	err = i.Run(context.Background(), s)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Have run scale function, with output status %s\n", s.Context.Status)

	err = os.WriteFile("output.jpeg", s.Context.Frame, 0644)
	if err != nil {
		panic(err)
	}

	for _, d := range s.Context.Detections {
		fmt.Printf("Detection at (%d,%d) scale %d\n", d.Col, d.Row, d.Scale)
	}
}
