package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"strconv"
	"time"
)

type akcii struct {
	name string
	fare float64
	date time.Time
}

type outp struct {
	start float64
	end   float64
	minzn float64
	maxzn float64
	fl    bool
}

type tick struct {
	date   time.Time
	tikers map[string]outp
}

func addToOutp(ak akcii, zap outp) outp {
	if zap.fl {
		if ak.fare > zap.maxzn {
			zap.maxzn = ak.fare
		}
		if ak.fare < zap.minzn {
			zap.minzn = ak.fare
		}
		zap.end = ak.fare
	} else {
		zap.end = ak.fare
		zap.start = ak.fare
		zap.maxzn = ak.fare
		zap.minzn = ak.fare
		zap.fl = true
	}
	return zap
}

func checkTick(sr tick) bool {
	for _, el := range sr.tikers {
		if el.fl {
			return true
		}
	}
	return false
}

func clearTick(t tick) tick {
	t.tikers = make(map[string]outp)
	var tim time.Time
	t.date = tim
	return t
}

func working3(el akcii, delta time.Duration, temp tick, out chan tick, timeP time.Time) (tick) {
	timeP = time.Date(el.date.Year(), el.date.Month(), el.date.Day(), 7, 0, 0, 0, time.UTC)

	if el.date.Before(timeP) {
		if checkTick(temp) {
			out <- temp
			temp = clearTick(temp)
		}

		return temp
	}

	for timeP.Add(delta).Before(el.date) {
		timeP = timeP.Add(delta)

	}

	if !((temp.date.After(timeP) || temp.date == timeP) && (temp.date.Before(timeP.Add(delta)) || temp.date == timeP.Add(delta))) {
		if checkTick(temp) {
			out <- temp
			temp = clearTick(temp)
		}
		temp.date = timeP
	}

	temp.tikers[el.name] = addToOutp(el, temp.tikers[el.name])

	return temp
}

func stage1(name string, done <-chan struct{}) <-chan akcii {

	out := make(chan akcii)

	csvFile, err := os.Open(name)
	if err != nil {
		fmt.Println(err, "problem with opening file")
	}

	r := csv.NewReader(bufio.NewReader(csvFile))
	var (
		akc  akcii
		fare float64
		ti   time.Time
	)
	layout := "2006-01-02 15:04:05"

	go func() {

		defer close(out)

		for {

			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err, "problem with read csv")
			}

			fare,err = strconv.ParseFloat(record[1], 64)
			if err != nil {
				fmt.Println(err, "problem with float parsing")
			}

			ti, err = time.Parse(layout, record[3])
			if err != nil {
				fmt.Println(err, "problem with time parsing")
			}

			akc = akcii{
				name: record[0],
				fare: fare,
				date: ti,
			}

			out <- akc

			select {
			case <-done:
				return
			default:
			}
		}
	}()

	return out
}

func stage2(deltas []time.Duration, done <-chan struct{}, in <-chan akcii) (<-chan tick, <-chan tick, <-chan tick) {

	out5 := make(chan tick)
	out30 := make(chan tick)
	out240 := make(chan tick)

	var timeSt time.Time

	temp := make([]tick, 3)
	go func() {

		for ii := 0; ii < 3; ii++ {
			temp[ii] = clearTick(temp[ii])
		}

		defer close(out5)
		defer close(out30)
		defer close(out240)

		for akc := range in {

			temp[0] = working3(akc, deltas[0], temp[0], out5, timeSt)
			temp[1] = working3(akc, deltas[1], temp[1], out30, timeSt)
			temp[2] = working3(akc, deltas[2], temp[2], out240, timeSt)

			select {
			case <-done:
				return
			default:
			}
		}
		out5 <- temp[0]
		out30 <- temp[1]
		out240 <- temp[2]

	}()

	return out5,out30,out240

}

func stage3(done chan struct{}, done2 chan struct{}, in ...<-chan tick) error {

	namesCandles := make([]string, 3)
	csvfile := make([]*os.File, 3)
	csvWriter := make([]*csv.Writer, 3)

	var err error

	namesCandles[0] = "hw3/candles_5min.csv"
	namesCandles[1] = "hw3/candles_30min.csv"
	namesCandles[2] = "hw3/candles_240min.csv"

	for i := 0; i < 3; i++ {
		csvfile[i], err = os.Create(namesCandles[i])
		if err != nil {
			return errors.Wrapf(err, "can't create file %s", namesCandles[i])
		}
		csvWriter[i] = csv.NewWriter(csvfile[i])
	}

	for k := 0; k < 3; k++ {

		go func(j int) {

			num := j
			defer csvfile[num].Close()
			for akc := range in[num] {
				for i, el := range (akc.tikers) {
					record := []string{i, akc.date.Format("2006-01-02T15:04:05Z"),
						strconv.FormatFloat(el.start, 'f', -1, 64),
						strconv.FormatFloat(el.maxzn, 'f', -1, 64),
						strconv.FormatFloat(el.minzn, 'f', -1, 64),
						strconv.FormatFloat(el.end, 'f', -1, 64)}

					err := csvWriter[num].Write(record)
					if err != nil {
						fmt.Printf("problem with writing to csv num=%d, record=%s", num, record)
					}
				}
				csvWriter[num].Flush()

				select {
				case <-done:
					return
				default:
				}
			}
			done2 <- struct{}{}
		}(k)
	}

	return nil

}

func testEnd(done2 chan struct{}, done chan struct{}) {
	kol := 3
	for {
		select {
		case <-done2:
			kol--
			if kol == 0 {
				done <- struct{}{}
				return
			}

		case <-done:
			return
		}
	}

}

func main() {

	var name string
	flag.StringVar(&name, "file", "hw3/trades.csv", "The name of news")
	flag.Parse()

	done2 := make(chan struct{})

	done := make(chan struct{})
	defer close(done)

	deltas := make([]time.Duration, 3)
	deltas[0] = 5 * time.Minute
	deltas[1] = 30 * time.Minute
	deltas[2] = 240 * time.Minute



	outst1:= stage1(name, done)

	out5, out30,out240 := stage2(deltas, done, outst1)

	err := stage3(done, done2,out5, out30,out240 )
	if err != nil {
		fmt.Println(err)
	}

	go testEnd(done2, done)

	//t0 := time.Now();
	select {
	case <-time.After(5 * time.Second):
		done <- struct{}{}
		//fmt.Println("Time out")
	case <-done:
	}
	//fmt.Printf("Elapsed time: %v\n", time.Since(t0))

}
