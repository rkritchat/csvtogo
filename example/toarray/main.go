package main

import (
	"fmt"
	"github.com/rkritchat/csvtogo"
	"log"
)

type CustInfo struct {
	Firstname string `json:"firstname" max:"10" min:"1"`
	Lastname  string `json:"lastname" min:"1"`
	Age       int    `json:"age" min:"1" max:"3"`
}

func main() {
	c, err := csvtogo.NewClient[CustInfo](
		"./test.csv",
		&csvtogo.Options{
			SkipHeader: true,
			ChunkSize:  1,
			Comma:      ',',
			SkipCol: []int{ //skip convert to struct on column 0 and 3
				0,
				3,
			},
		})
	if err != nil {
		log.Fatalln(err)
	}

	//simple convert csv to []T
	r, err := c.CsvToStruct()
	if err != nil {
		fmt.Println(err)
	}
	for _, val := range r {
		fmt.Println(val)
	}
}
