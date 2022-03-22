package main

import (
	"log"
	"os"

	pb "github.com/clintjedwards/gofer/proto"
	"google.golang.org/protobuf/proto"
)

// environment variables will be separated based on where they came from
// GOFER_NOTIFIER_USER_CONFIG_<NAME>
// GOFER_NOTIFIER_MAIN_CONFIG_<NAME>
// Notifiers will also get two special environment variables unique to them
// GOFER_NOTIFIER_RUN_REPORT=
// GOFER_API_TOKEN=
// The notifier run report is a byte string that can be turned into a proto.RunReport

// getCurrentRun converts a json string of the current run into a run object.
func getRunReport(report string) (*pb.RunReport, error) {
	currentRun := &pb.RunReport{}
	err := proto.Unmarshal([]byte(report), currentRun)
	if err != nil {
		return nil, err
	}

	return currentRun, nil
}

func main() {
	reportString := os.Getenv("GOFER_NOTIFIER_RUN_REPORT")
	report, err := getRunReport(reportString)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("notification successful for run %d\n", report.Id)
}
