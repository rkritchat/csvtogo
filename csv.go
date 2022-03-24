package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type CsvToStruct interface {
	Execute2() error
}

type csvToStruct[T any] struct {
	targetCSV string
	outChan   chan []T
	end       chan bool
	next      chan bool
}

func NewCsvToStruct[T any](outChan chan []T) (CsvToStruct, chan bool, chan bool) {
	end := make(chan bool, 1)
	next := make(chan bool, 1)
	return &csvToStruct[T]{
		outChan: outChan,
		end:     end,
		next:    next,
	}, end, next
}

func (c *csvToStruct[T]) Execute2() error {
	data := [][]string{
		{"AAAA0", "BBBB0", "111"},
		{"CCCC1", "DDDD1", "222"},
		{"CCCC2", "DDDD2", "333"},
		{"CCCC3", "DDDD3", "444"},
		{"CCCC4", "DDDD4", "555"},
		{"CCCC5", "DDDD5", "666"},
		{"CCCC6", "DDDD6", "999"},
	}
	batchSize := 1
	ref := make([]T, 1)
	var out []T
	for i, val := range data {
		tmp := ref[0] //copy reference
		for j, _ := range val {
			f := reflect.ValueOf(&tmp).Elem().Field(j)
			switch f.Interface().(type) {
			case int:
				v, err := strconv.Atoi(data[i][j])
				if err != nil {
					fmt.Println("invalid csv type, recording to struct it should be type int")
					return err
				}
				f.SetInt(int64(v))
			case string:
				f.SetString(data[i][j])
			case bool:
				b, err := strconv.ParseBool(data[i][j])
				if err != nil {
					fmt.Println("invalid bool type, recording to struct it should be type bool")
					return err
				}
				f.SetBool(b)
			case float64:
				v, err := strconv.ParseFloat(data[i][j], 64)
				if err != nil {
					fmt.Println("invalid bool type, recording to struct it should be type bool")
					return err
				}
				f.SetFloat(v)
			case float32:
				v, err := strconv.ParseFloat(data[i][j], 32)
				if err != nil {
					fmt.Println("invalid bool type, recording to struct it should be type bool")
					return err
				}
				f.SetFloat(v)
			default:
				fmt.Println("not found")
			}
		}
		out = append(out, tmp)
		if len(out) == batchSize {
			c.outChan <- out
			<-c.next //w8 until client ready to move
			out = []T{}
		}
	}
	if len(out) > 0 {
		c.outChan <- out
		<-c.next //w8 until client ready to move
	}

	c.end <- true
	return nil
}

//func (c *csvToStruct[T]) Execute() error {
//	f, err := os.Open(c.targetCSV)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	v := reflect.ValueOf(c.out)
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
//
//		reflect.ValueOf(&c.out).Elem().Field(0).SetString("1234")
//
//	}
//	return nil
//}
