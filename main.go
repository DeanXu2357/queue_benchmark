package main

import (
	"time"

	"github.com/Arafatk/glot"
)

var queueSize int
var numMessages int
var numConsumers int = 1

func main() {
	p, errP := glot.NewPlot(3, false, false)
	if errP != nil {
		panic(errP)
	}

	if err := p.SetTitle("Example Plot"); err != nil {
		panic(err)
	}

	if err := p.SetXrange(0, 100000); err != nil {
		panic(err)
	}
	if err := p.SetXLabel("Number of messages"); err != nil {
		panic(err)
	}
	if err := p.SetYLabel("time"); err != nil {
		panic(err)
	}

	for i := 0; i < 100000; i++ {
		queueSize = i

		// execute 100 times
		exeTime := make([]int64, 0, 100)
		for j := 0; j < 100; j++ {
			t := time.Now()
			numMessages = getNumberOfMessages()

			d := time.Since(t)
			exeTime = append(exeTime, d.Nanoseconds())
		}

		// exeTime 去除極端值＆平均

	}

	n := make([]float64, 1000)
	if err := p.AddPointGroup("Sample 1", "points", n); err != nil {
		panic(err)
	}
	if err := p.SetFormat("png"); err != nil {
		panic(err)
	}
	if err := p.SavePlot("gnuplot.png"); err != nil {
		panic(err)
	}

	if err := p.Close(); err != nil {
		panic(err)
	}
}

func getNumberOfMessages() int {
	return queueSize * 2
}
