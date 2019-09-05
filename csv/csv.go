package csv

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func ParseCsv(filename string) (c chan []string) {
	c = make(chan []string)

	log.Printf("Reading %s", filename)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.ReuseRecord = true

	for {
		record, err := reader.Read()
		if (err == io.EOF) {
			close(c)
			break
		}
		if (err != nil) {
			log.Fatal(err)
		}

		c <- record
	}

	log.Printf("Finished reading %s", filename)

	return
}
//
//func PutReport(f string, trades chan model.ReportedTrade) {
//	fixFilename := getFixFilename(f)
//
//	log.Debugf("Writing %s", fixFilename)
//	file, err := os.Create(fixFilename)
//	if err != nil {
//		log.Fatalf("cant write %s %s", f, err.Error())
//	}
//
//	defer file.Close()
//
//	_, err = file.WriteString("Deal,login,\"transaction time\",type,symbol,volume,open_price,close_price,profit,login_name,lei\n")
//	if err != nil {
//		log.Fatalf("cant write %s %s", fixFilename, err.Error())
//	}
//	for t := range trades {
//		//log.Debugf("append to %s", fixFilename)
//		_, err = file.WriteString(t.ToCSVString() + "\n")
//		if err != nil {
//			log.Fatalf("cant write %s %s", fixFilename, err.Error())
//		}
//	}
//}
//
//func getFixFilename(path string) string {
//	return filepath.Dir(path) + string(filepath.Separator) + FixFilePrefix + filepath.Base(path)
//}
