package csvtogo

import (
	"encoding/csv"
	"io"
	"os"
	"sync"
)

func csvReader[T any](csvFile string, comma rune, valueSetter func(T, []string, int) error) error {
	f, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer f.Close()

	var d []string
	var wg sync.WaitGroup
	var dChan = make(chan bool, 1)
	var pool = make(chan bool, 10)
	var chanErr = make(chan error, 1)
	defer close(pool)
	defer close(chanErr)

	reader := csv.NewReader(f)
	reader.Comma = comma
	row := -1
	ref := make([]T, 1)
	for {
		row += 1
		d, err = reader.Read()
		if err == io.EOF {
			//no more content
			break
		}
		pool <- true
		wg.Add(1)
		go asyncSet[T](valueSetter, ref[0], &wg, pool, chanErr, d, row)
	}

	go w8(&wg, &dChan)

	select {
	case <-dChan:
		return io.EOF
	case e := <-chanErr:
		return e
	}
}

func w8(wg *sync.WaitGroup, dChan *chan bool) {
	wg.Wait()
	close(*dChan)
}

func asyncSet[T any](valueSetter func(T, []string, int) error, ref T, wg *sync.WaitGroup, pool chan bool, chanErr chan error, d []string, row int) {
	defer func() {
		wg.Done()
		<-pool //release
	}()
	err := valueSetter(ref, d, row)
	if err != nil {
		chanErr <- err
	}
}
