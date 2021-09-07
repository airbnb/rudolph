package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"
)

// ParseCsvFile takes the full path to a .csv file and parses it. Each line yielded gets fed to an
// output channel that is expected to "do stuff" with it. The csv file MUST:
// * have a header row
// * be valid csv (consistent number of columns per row)
// * be newline delimited
func ParseCsvFile(
	filepath string,
	output chan<- map[string]string,
	wg *sync.WaitGroup,
) (err error) {
	f, err := os.Open(filepath)
	if err != nil {
		err = errors.Wrapf(err, "failed to open file: %s", filepath)
		return
	}
	defer f.Close() // this needs to be after the err check

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		err = errors.Wrapf(err, "failed to read csv file: %s", filepath)
		return
	}

	if len(lines) == 0 {
		err = errors.New("no lines detected in csv file")
		return
	}

	var header []string
	for lineNumber, line := range lines {
		if lineNumber == 0 {
			header = line
			continue
		}

		if len(line) != len(header) {
			err = errors.New(fmt.Sprintf("incorrect number of columns in line: %s", line))
			return
		}

		data := make(map[string]string, len(line))
		for columnNumber, item := range line {
			data[header[columnNumber]] = item
		}

		wg.Add(1)
		output <- data
	}

	return
}

func WriteCsvFile(
	filename string,
	headers []string,
	input <-chan []string,
	wg *sync.WaitGroup,
) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		err = errors.Wrapf(err, "failed to open file for writing: %s", filename)
		return
	}
	defer f.Close()
	err = f.Truncate(0)
	if err != nil {
		err = errors.Wrap(err, "file is not writable")
		return
	}

	csvWriter := csv.NewWriter(f)

	err = csvWriter.Write(headers)
	if err != nil {
		err = errors.Wrap(err, "could not write csv header")
		return
	}

	for {
		record := <-input

		err = csvWriter.Write(record)
		if err != nil {
			err = errors.Wrap(err, "could not write csv line")
			return
		}

		csvWriter.Flush()
		wg.Done()
	}
}
