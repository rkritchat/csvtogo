package main

import (
	"fmt"
	"github.com/rkritchat/csvtogo"
	"io"
	"log"
	"time"
)

type CustInfo struct {
	Firstname string `json:"firstname" max:"10" min:"1"`
	Lastname  string `json:"lastname" min:"0"`
	Age       int    `json:"age" min:"1" max:"3"`
	Married   bool   `json:"married"`
}

func main() {
	c, err := csvtogo.NewClient[CustInfo](
		"./test.csv",
		&csvtogo.Options{
			SkipHeader: true,
			Comma:      ',',
			SkipCols: []int{ //skip convert to struct on column 0 and 3
				0,
				3,
			},
		})
	if err != nil {
		log.Fatalln(err)
	}

	//convert csv row by row
	rows := c.CsvToRows()
	var r []*CustInfo
	defer rows.Close()
	for rows.Next() {
		tmp, err := rows.Read()
		if err == io.EOF { //return EOF is no more row
			fmt.Println("EOF")
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		if tmp != nil {
			//you can adjust struct here if needed
			fmt.Printf("%#v\n", tmp)
			fmt.Println("process something 2 secs")
			time.Sleep(2 * time.Second)
			r = append(r, tmp)
		}
	}

	//work with your lovely struct here
	fmt.Println(r)
}
