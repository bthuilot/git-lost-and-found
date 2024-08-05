package reporter

import (
	"encoding/json"
	"os"

	"github.com/bthuilot/git-scanner/pkg/processor"
	"github.com/sirupsen/logrus"
)

func groupBySecrets(results []processor.GitleaksResult, outputPath string) map[string][]processor.SecretsReport {
	secretsMap := make(map[string][]processor.SecretsReport)

	for _, report := range results {
		if _, ok := secretsMap[report.Match]; !ok {
			secretsMap[report.Match] = make([]processor.SecretsReport, 0)
		}
		secretsReport := processor.SecretsReport{}
		secretsReport.BlobHash = report.Fingerprint
		secretsReport.FileName = report.File
		file, err := os.ReadFile(outputPath + "/rawblobs/" + report.File)
		if err != nil {
			logrus.Error(err)
			file = []byte{}
		}
		secretsReport.RawBlob = string(file)
		secretsReport.Results = append(secretsReport.Results, report)

		secretsMap[report.Match] = append(secretsMap[report.Match], secretsReport)
	}

	return secretsMap
}

func WriteReport(output *os.File, results []processor.GitleaksResult, outputPath string) error {
	secretsMap := groupBySecrets(results, outputPath)
	logrus.Infof("Raw report %#v", results)
	logrus.Infof("Writing report with %d commits", len(secretsMap))
	logrus.Infof("Writing to file %s", output.Name())
	jsonBytes, err := json.MarshalIndent(secretsMap, "", "  ")
	if err != nil {
		return err
	}
	_, err = output.Write(jsonBytes)
	if err != nil {
		return err
	}
	return nil
}
