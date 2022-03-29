package main

import (
	"fmt"
	"github.com/rkritchat/csvtogo"
	"log"
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
	Firstname string `json:"firstname" max:"10" min:"1"` //<-- put max / min here if you need.
	Lastname  string `json:"lastname" min:"1"`
	Age       int    `json:"age" min:"1" max:"3"`
	Married   bool   `json:"married"`
}

func main() {
	c, err := csvtogo.NewClient[CustInfo](
		"./test.csv",
		&csvtogo.Options{ //the option is an optional, csvtogo will use default if ops is nil
			SkipHeader: true,
			Comma:      ',',
			SkipCols: []int{ //skip convert to struct at column 0 and 3
				0, //SKIP COLUMN ID
				3, //SKIP COLUMN ADDR
			},
		})
	if err != nil {
		log.Fatalln(err)
	}

	//simple convert csv to []T, if you need to adjust the value for some rows, suggest to use csv to rows instead
	r, err := c.CsvToStruct()
	if err != nil {
		fmt.Println(err)
	}
	for _, val := range r {
		fmt.Println(val)
	}
}
