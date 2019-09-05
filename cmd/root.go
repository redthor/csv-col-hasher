package cmd

import (
	"github.com/redthor/csv-col-hasher/csv"
	"github.com/spf13/cobra"
	"log"
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func hash(cmd *cobra.Command, args []string) {
	c := csv.ParseCsv(csvFile)

	// Check that the file has some records
	select {
		case first := <- c:
			if (int(colNum) > len(first)) {
				log.Fatalf("Col number %d is out of range of record length %d", colNum, len(first))
			}
		default:
			log.Fatal("No records")
	}

	for csvLine := range csv.ParseCsv(csvFile) {
		log.Println(csvLine[colNum])
	}
}