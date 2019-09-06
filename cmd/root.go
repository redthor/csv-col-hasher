package cmd

import (
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "csv-col-hasher",
	Short: "Convert a column in a CSV file to a hashed value",
	Long:  "",
	Run:   hash,
}

var csvFile string
var outputFile string
var colNum uint16

func init() {
	cobra.OnInitialize()
	rootCmd.Flags().StringVarP(&csvFile, "csv-file", "f", "", "The CSV file to read")
	_ = rootCmd.MarkFlagRequired("csv-file")
	rootCmd.Flags().Uint16VarP(&colNum, "col-num", "n", 0, "The column to replace with hashed values, first is 0")
	_ = rootCmd.MarkFlagRequired("col-num")
	rootCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "The CSV file to write to. If not set it will print to stdout")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func hash(cmd *cobra.Command, args []string) {
	log.Printf("Processing column %d in '%s'", colNum, csvFile)

	out := createOutputWriter(outputFile)
	defer out.Flush()
	c := make(chan []string)
	go parseCsv(csvFile, c)

	h := sha1.New()
	first := true

	for csvLine := range c {
		if first {
			if err := out.Write(csvLine); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			first = false
			continue
		}
		toConvert := csvLine[colNum]
		h.Write([]byte(toConvert))
		csvLine[colNum] = hex.EncodeToString(h.Sum(nil))
		if err := out.Write(csvLine); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}

		h.Reset()
	}
	out.Flush()

	if err := out.Error(); err != nil {
		log.Fatalf("Error flushing csv '%s'", err.Error())
	}
}

func createOutputWriter(filename string) *csv.Writer {
	// If not output file, use stdOut
	if outputFile == "" {
		return csv.NewWriter(os.Stdout)
	}

	out, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error opening output file '%s'", outputFile)
	}

	log.Printf("Writing to '%s'", filename)

	return csv.NewWriter(out)
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
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		c <- record
	}
}
