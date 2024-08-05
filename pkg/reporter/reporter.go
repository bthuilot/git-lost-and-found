package reporter

import (
	"encoding/json"
	"os"

	"github.com/bthuilot/git-scanner/pkg/processor"
	"github.com/sirupsen/logrus"
)

func groupBySecrets(results []processor.GitleaksResult) map[string][]processor.SecretsReport {
	secretsMap := make(map[string][]processor.SecretsReport)

	for _, report := range results {
		if _, ok := secretsMap[report.Match]; !ok {
			secretsMap[report.Match] = make([]processor.SecretsReport, 0)
		}
		secretsReport := processor.SecretsReport{}
		secretsReport.BlobHash = report.Fingerprint
		secretsReport.FileName = report.File
		secretsReport.RawBlob = report.RawBlob
		secretsReport.Results = append(secretsReport.Results, report)

		secretsMap[report.Match] = append(secretsMap[report.Match], secretsReport)
	}

	return secretsMap
}

func WriteReport(output *os.File, results []processor.GitleaksResult) error {
	secretsMap := groupBySecrets(results)
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
