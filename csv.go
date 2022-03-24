package main

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)

type CsvToStruct[T any] struct {
	targetCSV string
	chunkSize int
	outChan   chan []T
	end       chan bool
	next      chan bool
	run       bool
}

func Conv[T any](csvFile string) CsvToStruct[T] {
	c := CsvToStruct[T]{
		targetCSV: csvFile,
		outChan:   make(chan []T, 1),
		chunkSize: 2,
		end:       make(chan bool, 1),
		next:      make(chan bool, 1),
		run:       true,
	}
	c.start()
	return c
}

func (c *CsvToStruct[T]) Next() bool {
	if c.run {
		c.next <- true
		return true
	}

	return true
}

func (c *CsvToStruct[T]) Read() ([]T, error) {
	for c.run {
		select {
		case data := <-c.outChan:
			return data, nil
		case <-c.end:
			fmt.Println("no more data, stop..")
			c.run = false
			return nil, io.EOF
		default:
			//waiting
		}
	}
	return nil, nil
}

func (c *CsvToStruct[T]) start() {
	go func() {
		//TODO run with goroutine here and remove return err, if error occurred then send it to chan
		data := [][]string{
			{"AAAA0", "BBBB0", "111"},
			{"CCCC1", "DDDD1", "222"},
			{"CCCC2", "DDDD2", "333"},
			{"CCCC3", "DDDD3", "444"},
			{"CCCC4", "DDDD4", "555"},
			{"CCCC5", "DDDD5", "666"},
			{"CCCC6", "DDDD6", "999"},
		}
		ref := make([]T, 1)
		var out []T
		for i, val := range data {
			tmp := ref[0] //copy reference
			for j := range val {
				f := reflect.ValueOf(&tmp).Elem().Field(j)
				err := set(f, data[i][j])
				if err != nil {
					log.Fatalln(err) //TODO fix me by sending err to chan
					//return err
				}
			}
			out = append(out, tmp)
			c.send(&out)
		}

		//last chunk
		c.send(&out, len(out) > 0)
		c.end <- true
	}()
}

func (c *CsvToStruct[T]) send(out *[]T, force ...bool) {
	if c.chunkSize == len(*out) || (len(force) > 0 && force[0]) {
		c.outChan <- *out
		<-c.next //w8 until client ready to move
		*out = []T{}
	}
}

func set(f reflect.Value, val string) error {
	switch f.Interface().(type) {
	case int, int32, int64:
		v, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("invalid csv type, recording to struct it should be type int")
			return err
		}
		f.SetInt(int64(v))
	case string:
		f.SetString(val)
	case bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			fmt.Println("invalid bool type, recording to struct it should be type bool")
			return err
		}
		f.SetBool(b)
	case float32, float64:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Println("invalid bool type, recording to struct it should be type bool")
			return err
		}
		f.SetFloat(v)
	default:
		fmt.Println("not found")
	}
	return nil
}

func (c *CsvToStruct[T]) Close() {
	close(c.next)
	close(c.outChan)
	close(c.end)
}

//func (c *CsvToStruct[T]) Execute() error {
//	f, err := os.Open(c.targetCSV)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	reader := csv.NewReader(f)
//	counter := 0
//	for {
//		tmp, err := reader.Read()
//		if err == io.EOF {
//			break
//		}
//		counter += 1
//		if len(tmp) != v.NumField() {
//			fmt.Println(counter)
//			return fmt.Errorf("invalid at %v", counter)
//		}
//	}
//	return nil
//}
