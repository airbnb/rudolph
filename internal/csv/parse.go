package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
)

// ParseCsvFile takes the full path to a .csv file and parses it. Each line yielded gets fed to an
// output channel that is expected to "do stuff" with it. The csv file MUST:
// * have a header row
// * be valid csv (consistent number of columns per row)
// * be newline delimited
func ParseCsvFile(filepath string) (output chan map[string]string, err error) {
	output = make(chan map[string]string)

	f, err := os.Open(filepath)
	if err != nil {
		err = errors.Wrapf(err, "failed to open file: %s", filepath)
		return
	}

	reader := csv.NewReader(f)
	header, err := reader.Read()
	if err != nil {
		err = errors.Wrapf(err, "failed to read csv file: %s", filepath)
		return
	}

	go func() {
		defer f.Close() // this needs to be done here, or will result in error reading a closed file
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break // reached EOF
			} else if err != nil {
				panic(err) // This is likely an unrecoverable error, so just quit
			}

			if len(line) != len(header) {
				// Log the error. consider using a struct with err atttribute, or
				// list of errors, to track overall errors
				fmt.Println("[ERROR] incorrect number of columns in line:", line)
				continue
			}
			data := make(map[string]string, len(line))
			for columnNumber, item := range line {
				data[header[columnNumber]] = item
			}

			output <- data
		}
		close(output) // Calling close here ensures workers exit properly
	}()

	return
}

func WriteCsvFile(
	filename string,
	headers []string,
	input <-chan []string,
) (wg *sync.WaitGroup, err error) {
	f, err := os.Create(filename)
	if err != nil {
		err = errors.Wrapf(err, "failed to open file for writing: %s", filename)
		return
	}
	err = f.Truncate(0)
	if err != nil {
		f.Close()
		err = errors.Wrap(err, "file is not writable")
		return
	}

	csvWriter := csv.NewWriter(f)

	err = csvWriter.Write(headers)
	if err != nil {
		err = errors.Wrap(err, "could not write csv header")
		return
	}

	wg = new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer f.Close()
		defer csvWriter.Flush()

		for record := range input {
			err = csvWriter.Write(record)
			if err != nil {
				err = errors.Wrap(err, "could not write csv line")
				return
			}
		}
	}()

	return
}
