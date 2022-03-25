package csvtogo

import (
	"encoding/csv"
	"io"
	"os"
)

func csvReader[T any](csvFile string, comma rune, valueSetter func(r T, d []string, c int) error) error {
	f, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer f.Close()

	var d []string
	reader := csv.NewReader(f)
	reader.Comma = comma
	counter := -1
	ref := make([]T, 1)
	for {
		counter += 1
		d, err = reader.Read()
		if err == io.EOF {
			//no more content
			break
		}

		//setter logic
		err = valueSetter(ref[0], d, counter)
		if err != nil {
			return err
		}
	}
	return nil
}

func readTxt() {
	//TODO implement me
}
