package main

import (
	"fmt"
	"github.com/rkritchat/csvtogo"
	"io"
	"log"
	"time"
)

//the example test.csv file
//  +--------------------------------------------------------------------+
//  |  ID  ,  FIRSTNAME  ,  LASTNAME   ,  ADDR       ,  AGE  ,  MARRIED  |   <- header
//  |  1   ,  Kritchat   ,  Rojanaphruk,  test-addr1 ,  2    ,  false    |
//  |  2   ,  Uefa       ,  Uef        ,  test-addr2 ,  1    ,  true     |
//  |  3   ,  Luffy      ,  Luf        ,  test-addr3 ,  21   ,  false    |
//  +--------------------------------------------------------------------+
// PS. comma can set to any rune type such as pipe ( | ) or something you wish. (Options.Comma: '|')

//CustInfo
//convert above csv file to CustInfo struct order by column and struct field
//Ps. Struct or field name is no need to match which column
//Ps2. csvtogo support validating value such as max length / min length.
type CustInfo struct {
	Name     string    `json:"name" max:"10" min:"1"` //<-- put max / min here if you need.
	Lastname string    `json:"lastname" min:"0"`
	Age      int       `json:"age" min:"1" max:"3"`
	Married  bool      `json:"married"`
	CreateAt time.Time `json:"create_at"` //You can add addition field here also.
}

func main() {
	c, err := csvtogo.NewClient[CustInfo](
		"./test.csv",
		&csvtogo.Options{ //the option is an optional, csvtogo will use default if ops is nil
			SkipHeader: true, //SKIP header
			Comma:      ',',
			SkipCols: []int{ //skip convert to struct at column 0 and 3
				0, //SKIP COLUMN ID
				3, //SKIP COLUMN ADDR
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
