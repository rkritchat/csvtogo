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
	rows := Conv[CustInfo]()
	for rows.Next() { //next call <- next
		tmp, err := rows.Read() //return EOF is no more rows, return []T, err
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		if tmp != nil {
			fmt.Printf("%#v\n", tmp)
			fmt.Println("process something 2 secs")
			time.Sleep(1 * time.Second)
		}
	}
}
