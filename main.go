package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Arafatk/glot"

	"qbench/contract"
	"qbench/implbych"
	"qbench/implbyrb"
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
	if err := p.SetYLabel("time(Nanoseconds)"); err != nil {
		panic(err)
	}

	ctx := context.Background()

	n := make([]float64, 1000)
	for i := 0; i < 20000; i++ {
		fmt.Println("ch impl", i)

		// execute 100 times
		exeTime := make([]int64, 0, 100)
		for j := 0; j < 100; j++ {
			fmt.Println("ch impl", i, j)
			q := implbych.NewChImpl(int64(i))

			s, errS := testQueue(ctx, q, i, 1)
			if errS != nil {
				panic(errS)
			}

			exeTime = append(exeTime, s)
			fmt.Println("time:", s)
		}

		// 執行時間去除極端值＆平均
		e := removeMaxMinAndAverage(exeTime)
		n[i] = float64(e)
	}
	if err := p.AddPointGroup("Impl by channel", "points", n); err != nil {
		panic(err)
	}

	n1 := make([]float64, 1000)
	for i := 0; i < 20000; i++ {
		fmt.Println("rb impl", i)

		// execute 100 times
		exeTime := make([]int64, 0, 100)
		for j := 0; j < 100; j++ {
			q := implbyrb.NewRingBufferImpl(int64(i))

			s, errS := testQueue(ctx, q, i, 1)
			if errS != nil {
				panic(errS)
			}

			exeTime = append(exeTime, s)
		}

		// 執行時間去除極端值＆平均
		e := removeMaxMinAndAverage(exeTime)
		n1[i] = float64(e)
	}
	if err := p.AddPointGroup("Impl by ring buffer", "points", n1); err != nil {
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

func removeMaxMinAndAverage(exeTime []int64) int64 {
	if len(exeTime) <= 3 {
		return averageFloat(exeTime)
	}

	maxF := exeTime[0]
	minF := exeTime[0]
	for i := range exeTime {
		if exeTime[i] > maxF {
			maxF = exeTime[i]
		}
		if exeTime[i] < minF {
			minF = exeTime[i]
		}
	}
	var fs []int64
	for i := range exeTime {
		if exeTime[i] != maxF && exeTime[i] != minF {
			fs = append(fs, exeTime[i])
		}
	}

	return averageFloat(fs)
}

func averageFloat(s []int64) int64 {
	if len(s) == 0 {
		return 0
	}

	var sum int64
	for i := range s {
		sum += s[i]
	}
	return sum / int64(len(s))
}

func testQueue(ctx context.Context, q contract.Queue, bufferSize, numOfConsumer int) (int64, error) {
	wg := sync.WaitGroup{}

	now := time.Now()
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < bufferSize*numOfConsumer; i++ {
			if err := q.Put(ctx, strconv.Itoa(i)); err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	rtn := make([]int, numOfConsumer)
	for i := range rtn {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < bufferSize; j++ {
				if err := q.Consume(ctx, func(s []string) error {
					rtn[i] += len(s)
					return nil
				}); err != nil {
					fmt.Println(err)
					return
				}

			}
		}()
	}

	wg.Wait()
	t := time.Since(now).Nanoseconds()

	var sum int
	for i := range rtn {
		sum += rtn[i]
	}
	if sum != bufferSize*numOfConsumer {
		return 0, fmt.Errorf("Expected %d messages, got %d", bufferSize, sum)
	}

	if err := q.Close(ctx); err != nil {
		return 0, fmt.Errorf("Error closing queue: %v", err)
	}

	return t, nil
}
