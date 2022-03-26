package csvtogo

import (
	"os"
	"unicode/utf8"
)

type Client[T any] struct {
	Executor[T]
}

func NewClient[T any](file string, ops ...*Options) (*Client[T], error) {
	option := _defaultOps
	if ops != nil {
		//validate ops
		options, err := validateOps(ops[0])
		if err != nil {
			return nil, err
		}
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

func validateOps(ops *Options) (*Options, error) {
	if ops == nil {
		return nil, nil
	}
	if utf8.RuneCountInString(string(ops.Comma)) == 0 {
		return nil, csvCommaIsRequired
	}
	if len(ops.SkipCol) > 0 {
		m := make(map[int]int)
		for _, val := range ops.SkipCol {
			m[val] = val
		}
		ops.skipper = m
	}
	return ops, nil
}
