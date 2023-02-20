package converter

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/SkobelevIgor/stackexchange-xml-converter/encoders"
)

func convertToCSV(typeName string, xmlFile *os.File, csvFileBasePath *os.File, cfg Config) (total int64, converted int64, err error) {
	csvWriter := nil // csv.NewWriter(csvFile)
	// defer csvWriter.Flush()

	encoder, err := encoders.NewEncoder(typeName)
	if err != nil {
		return
	}

	err = csvWriter.Write(encoder.GetCSVHeaderRow())
	if err != nil {
		return
	}

	iterator := NewIterator(xmlFile)

	var iErr error
	for iterator.Next() {
		total++

		if total % cfg.BatchSize == 1) {
			if csvWRiter != nil {
				csvWriter.Flush()
			}
			batchNum := total / cfg.BatchSize
			resultFilePath := fmt.Sprintf("%s-%06d", csvFileBasePath, batchNum)
			csvWriter = csv.NewWriter(resultFilePath)
		}

		encoder, _ := encoders.NewEncoder(typeName)
		iErr = iterator.Decode(&encoder)
		if iErr != nil {
			log.Printf("[%s] Error: %s", typeName, iErr)
			continue
		}

		if cfg.SkipHTMLDecoding {
			encoder.EscapeFields()
		}

		iErr = csvWriter.Write(encoder.GETCSVRow())
		if iErr != nil {
			log.Printf("[%s] Error: %s", typeName, iErr)
			continue
		}
		converted++
	}

	return
}
