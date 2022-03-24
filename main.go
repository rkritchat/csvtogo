package main

import (
	"fmt"
	"log"
	"time"
)

type CustInfo struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

func main() {
	//c := CustInfo{}
	//reflect.ValueOf(&c).Elem().Field(0).SetString("1234")
	//v := reflect.ValueOf(c)
	//typeOfS := v.Type()
	//fmt.Println(v.NumField())
	//for i := 0; i < v.NumField(); i++ {
	//	fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	//}

	output := make(chan []CustInfo)
	defer close(output)

	c, end, next := NewCsvToStruct[CustInfo](output)
	go func() {
		err := c.Execute2()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	run := true
	for run {
		select {
		case data := <-output:
			fmt.Printf("%#v\n", data)
			fmt.Println("process something 2 secs")
			time.Sleep(2 * time.Second)
			next <- true
		case <-end:
			fmt.Println("no more data, stop..")
			run = false
		default:
			//waiting
		}
	}
}
