package cmd

import (
	"encoding/hex"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
	"encoding/csv"
	"crypto/sha1"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "csv-col-hasher",
	Short: "Convert a column in a CSV file to a hashed value",
	Long:  "",
	Run: hash,
}

var csvFile string
var colNum uint16

func init() {
	cobra.OnInitialize()
	rootCmd.Flags().StringVarP(&csvFile, "csv-file", "f", "", "The CSV file to read")
	_ = rootCmd.MarkFlagRequired("csv-file")
	rootCmd.Flags().Uint16VarP(&colNum, "col-num", "n", 0, "The column to replace with hashed values, first is 0")
	_ = rootCmd.MarkFlagRequired("col-num")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func hash(cmd *cobra.Command, args []string) {
	outputFile, err := ioutil.TempFile("", "output.*.csv")

	if (err != nil) {
		log.Fatal("Error opening output file")
	}
	defer outputFile.Close()

	log.Printf("Processing column %d in '%s' and writing to %s", colNum, csvFile, outputFile.Name())

	c := make(chan []string)
	go parseCsv(csvFile, c)

	h := sha1.New()
	writer := csv.NewWriter(outputFile)
	first := true

	for csvLine := range c {
		if (first) {
			if err := writer.Write(csvLine); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			first = false
			continue
		}
		toConvert := csvLine[colNum]
		h.Write([]byte(toConvert))
		csvLine[colNum] = hex.EncodeToString(h.Sum(nil))
		if err := writer.Write(csvLine); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}

		h.Reset()
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatalf("Error flushing csv '%s'", err.Error())
	}

	log.Printf("Finished writing '%s'", outputFile.Name())
}

func parseCsv(filename string, c chan []string) {
	defer close(c)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	for {
		record, err := reader.Read()
		if (err == io.EOF) {
			break
		}
		if (err != nil) {
			log.Fatal(err)
		}

		c <- record
	}
}