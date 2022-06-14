[![Coverage Status](https://coveralls.io/repos/github/rkritchat/csvtogo/badge.svg?branch=master)](https://coveralls.io/github/rkritchat/csvtogo?branch=master)

csvtogo used for convert csv file to struct by using reflex and generic, required Go version 1.18+

## Installation
```shell
go get -u github.com/rkritchat/csvtogo

```

## Usage
Convert csv to array.

```go
//the example sample.csv file
//  +--------------------------------------------------------------------+
//  |  ID  ,  FIRSTNAME  ,  LASTNAME   ,  ADDR       ,  AGE  ,  MARRIED  |   <- header
//  |  1   ,  John       ,  Doe        ,  test-addr1 ,  2    ,  false    |
//  |  2   ,  Uefa       ,  Uef        ,  test-addr2 ,  1    ,  true     |
//  |  3   ,  Luffy      ,  Luf        ,  test-addr3 ,  21   ,  false    |
//  +--------------------------------------------------------------------+
// PS. comma can set to any rune type such as pipe ( | ) or something related to csv file. (Options.Comma: '|')

//CustInfo
//convert above csv file to CustInfo struct order by column and struct field by using reflect and generic
//Ps. Struct or field name is no need to match which column
//Ps2. csvtogo support validating value such as max length / min length.
type CustInfo struct {
	Firstname string `json:"firstname" max:"10" min:"1"` //<-- put max / min tag here if needed.
	Lastname  string `json:"lastname" min:"1"`
	Age       int    `json:"age" min:"1" max:"3"`
	Married   bool   `json:"married"`
}

func main() {
	c, err := csvtogo.NewClient[CustInfo](
		"./sample.csv",
		&csvtogo.Options{ //the option is an optional, csvtogo will use default if ops is nil
			SkipHeader: true, //SKIP first row
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

```

Convert csv to struct row by row.

```go
func main() {
	c, err := csvtogo.NewClient[CustInfo](
		"./sample.csv",
		&csvtogo.Options{ //the option is an optional, csvtogo will use default if ops is nil
			SkipHeader: true, //SKIP first row
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
		if err != nil {
			if err == io.EOF { //return EOF is no more row
				fmt.Println("EOF")
				break
			}
			fmt.Println(err)
			break
		}

		if tmp != nil {
			//you can adjust struct here if needed
			fmt.Printf("%#v\n", tmp)
			fmt.Println("process something 1 secs")
			time.Sleep(1 * time.Second)
			r = append(r, tmp)
		}
	}

	//work with your lovely struct here
	fmt.Println(r)
}
```

MIT License

Copyright (c) 2022 rkritchat
