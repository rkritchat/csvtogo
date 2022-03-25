package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

const (
	defaultChunkSize = 1000
)

var _defaultOps = Options{
	SkipHeader: true,
}

type CsvToStruct[T any] struct {
	targetCSV string
	outChan   chan []T
	end       chan bool
	next      chan bool
	errChan   chan error
	run       bool
	ops       Options
}

type Options struct {
	ChunkSize   int
	SkipHeader  bool
	SkipCol     []string
	ReplaceWith map[string]string
}

func NewClient[T any](csvFile string, ops ...Options) CsvToStruct[T] {
	options := _defaultOps
	if ops != nil {
		options = ops[0]
	}
	return CsvToStruct[T]{
		targetCSV: csvFile,
		ops:       options,
		outChan:   make(chan []T, 1),
		end:       make(chan bool, 1),
		next:      make(chan bool, 1),
		errChan:   make(chan error),
		run:       true,
	}
}

func (c *CsvToStruct[T]) ToList() *CsvToStruct[T] {
	go c.start()
	return c
}

func (c *CsvToStruct[T]) Next() bool {
	if c.run {
		c.next <- true
		return true
	}
	return false
}

func (c *CsvToStruct[T]) Read() ([]T, error) {
	for c.run {
		select {
		case data := <-c.outChan:
			return data, nil
		case <-c.end:
			c.run = false
			return nil, io.EOF
		case err := <-c.errChan:
			c.run = false
			return nil, err
		default:
			//waiting
		}
	}
	return nil, nil
}

func (c *CsvToStruct[T]) start() {
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
	var isErr bool
	for _, val := range data {
		if isErr {
			break
		}
		tmp := ref[0] //copy reference
		err := c.setValue(val, &tmp)
		if err != nil {
			c.errChan <- err
		}
		out = append(out, tmp)
		c.send(&out)
	}
	if isErr {
		return
	}

	//last chunk
	c.send(&out, len(out) > 0)
	c.end <- true
}

func (c *CsvToStruct[T]) setValue(data []string, tmp *T) error {
	for i, val := range data {
		f := reflect.ValueOf(tmp).Elem().Field(i)
		err := typeSafe(f, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CsvToStruct[T]) send(out *[]T, force ...bool) {
	if len(*out) >= c.ops.ChunkSize || (len(force) > 0 && force[0]) {
		c.outChan <- *out
		<-c.next //w8 until client ready to move
		*out = []T{}
	}
}

func typeSafe(f reflect.Value, val string) error {
	switch f.Interface().(type) {
	case int, int32, int64:
		v, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("invalid csv type, recording to struct it should be type int") //TODO fixme, generate manual err
			return err
		}
		f.SetInt(int64(v))
	case string:
		f.SetString(val)
	case bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			fmt.Println("invalid bool type, recording to struct it should be type bool") //TODO fixme
			return err
		}
		f.SetBool(b)
	case float32, float64:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Println("invalid bool type, recording to struct it should be type bool") //TODO fixme
			return err
		}
		f.SetFloat(v)
	default:
		fmt.Println("not found") //TODO return err not support
	}
	return nil
}

func (c *CsvToStruct[T]) Close() {
	close(c.next)
	close(c.outChan)
	close(c.end)
}

func (c *CsvToStruct[T]) Rows() *CsvToStruct[T] {
	go func() {
		f, err := os.Open(c.targetCSV)
		if err != nil {
			c.errChan <- err
			return
		}
		defer f.Close()

		reader := csv.NewReader(f)
		counter := -1
		ref := make([]T, 1)
		var out []T
		for {
			counter += 1
			tmp, err := reader.Read()
			if err == io.EOF {
				break
			}
			//skip header
			if c.ops.SkipHeader && counter == 0 {
				continue
			}
			typeRef := ref[0]
			v := reflect.ValueOf(&typeRef).Elem()
			if len(tmp) != v.NumField() {
				fmt.Println(counter)
				c.errChan <- fmt.Errorf("invalid at %v", counter)
				break
			}
			err = c.setValue(tmp, &typeRef)
			if err != nil {
				c.errChan <- err
				break
			}
			out = append(out, typeRef)
			c.send(&out, true)
		}
		c.end <- true
	}()
	return c
}
