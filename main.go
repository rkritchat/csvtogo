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
	outChan := make(chan []CustInfo, 1)
	defer close(outChan)
	c := Conv[CustInfo](outChan)
	for c.Next() { //next call <- next
		tmp, err := c.Read() //return EOF is no more rows, return []T, err
		if err == io.EOF {
			fmt.Println("eofff")
			break
		}
		if tmp != nil {
			fmt.Printf("%#v\n", tmp)
			fmt.Println("process something 2 secs")
			time.Sleep(1 * time.Second)
		}
	}
}
