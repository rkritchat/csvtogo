package main

import (
	"fmt"
	"io"
	"time"
)

type CustInfo struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

func main() {
	client := NewClient[CustInfo](
		"./test.csv",
		Options{
			SkipHeader: true,
			ChunkSize:  1,
		})
	rows := client.Rows()
	defer rows.Close()
	for rows.Next() {
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
