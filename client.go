package csvtogo

import (
	"os"
)

type Client[T any] struct {
	Executor[T]
}

func NewClient[T any](file string, ops ...*Options) (*Client[T], error) {
	option := _defaultOps
	if ops != nil {
		options := ops[0]
		//convert SkipCols to map
		options.skipper = initSkipper(options.SkipCols)
		option = *options
	}

	//validate file
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	_ = f.Close()

	return &Client[T]{
		Executor[T]{
			file:     file,
			ops:      option,
			outsChan: make(chan []T, 1),
			outChan:  make(chan T, 1),
			endChan:  make(chan bool, 1),
			nextChan: make(chan bool, 1),
			errChan:  make(chan error),
			run:      true,
		},
	}, nil
}

func initSkipper(skipCols []int) map[int]int {
	if len(skipCols) > 0 {
		m := make(map[int]int)
		for _, val := range skipCols {
			m[val] = val
		}
		return m
	}
	return nil
}
