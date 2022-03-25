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

	rows := c.Rows()
	defer rows.Close()
	for rows.Next() { //TODO check if client break loop, it will error or not
		tmp, err := rows.Read() //return EOF is no more rows, return T, err
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		if tmp != nil {
			fmt.Printf("%#v\n", tmp)
			fmt.Println("process something 2 secs")
			time.Sleep(1 * time.Second)
		}
	}
}
