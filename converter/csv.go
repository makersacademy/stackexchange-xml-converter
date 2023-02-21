package converter

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/SkobelevIgor/stackexchange-xml-converter/encoders"
)

func getWriterForNewBatch(encoder encoders.Encoder, basePath string, batchNum int64) (resultFile *os.File, csvWriter *csv.Writer, err error) {
	resultFilePath := fmt.Sprintf("%s-%06d", basePath, batchNum)
	resultFile, err = os.Create(resultFilePath)
	if err != nil {
		return
	}
	csvWriter = csv.NewWriter(resultFile)

	err = csvWriter.Write(encoder.GetCSVHeaderRow())
	if err != nil {
		return
	}

	return
}

func convertToCSV(typeName string, xmlFile *os.File, csvFile *os.File, resultFileBasePath string, cfg Config) (total int64, converted int64, err error) {
	encoder, err := encoders.NewEncoder(typeName)
	if err != nil {
		return
	}

	resultFile, csvWriter, err := getWriterForNewBatch(encoder, resultFileBasePath, 1)
	if err != nil {
		return
	}

	iterator := NewIterator(xmlFile)

	var iErr error
	for iterator.Next() {
		total++

		if total > 1 && total%cfg.BatchSize == 1 {
			log.Printf("[%s] Starting new Batch: %d", typeName, total)
			csvWriter.Flush()
			resultFile.Close()

			batchNum := 1 + total/cfg.BatchSize
			resultFile, csvWriter, err = getWriterForNewBatch(encoder, resultFileBasePath, batchNum)
			if err != nil {
				return
			}
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

	csvWriter.Flush()
	resultFile.Close()

	return
}
