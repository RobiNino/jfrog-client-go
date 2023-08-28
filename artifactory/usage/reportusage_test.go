package usage

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsVersionCompatible(t *testing.T) {
	tests := []struct {
		artifactoryVersion string
		expectedResult     bool
	}{
		{"6.5.0", false},
		{"6.2.0", false},
		{"5.9.0", false},
		{"6.0.0", false},
		{"6.6.0", false},
		{"6.9.0", true},
		{utils.Development, true},
		{"6.10.2", true},
		{"6.15.2", true},
	}
	for _, test := range tests {
		t.Run(test.artifactoryVersion, func(t *testing.T) {
			err := utils.ValidateMinimumVersion(utils.Xray, test.artifactoryVersion, minArtifactoryVersion)
			if test.expectedResult {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, fmt.Sprintf(utils.MinimumVersionMsg, utils.Xray, test.artifactoryVersion, minArtifactoryVersion))
			}
		})
	}
}

func TestReportUsageJson(t *testing.T) {
	type reportUsageTestCase struct {
		productId      string
		commandName    string
		serviceId      ReportUsageAttribute
		serverSize     ReportUsageAttribute
		expectedResult string
		jsonPatternNum int
	}

	jsonPatterns := []string{
		`{"productId":"%s","features":[{"featureId":"%s","attributes":{"%s":"%s","%s":"%s"}}]}`,
		`{"productId":"%s","features":[{"featureId":"%s","attributes":{"%s":"%s"}}]}`,
		`{"productId":"%s","features":[{"featureId":"%s"}]}`,
	}

	preTests := []reportUsageTestCase{
		{"jfrog-cli-go/1.26.0", "rt_transfer_files", ReportUsageAttribute{"sourceServiceId", "jfrt@01g8dj3wcw22y01atqp63n1haq"}, ReportUsageAttribute{"sourceStorageSize", "6.08 GB"}, "{\"productId\":\"jfrog-cli-go/1.26.0\",\"features\":[{\"featureId\":\"rt_transfer_files\",\"attributes\":{\"sourceStorageSize\":\"6.08 GB\",\"sourceServiceId\":\"jfrt@01g8dj3wcw22y01atqp63n1haq\"}}]}", 0},
		{"jfrog-client-go", "rt_download", ReportUsageAttribute{}, ReportUsageAttribute{"sourceStorageSize", "3.58 GB"}, "{\"productId\":\"jfrog-client-go\",\"features\":[{\"featureId\":\"rt_download\"}]}", 1},
		{"test", "rt_build", ReportUsageAttribute{"sourceServiceId", "jfrt@01g8dj3wcw22y01atqp63n1haq"}, ReportUsageAttribute{}, "", 1},
		{"agent/1.25.0", "rt_go", ReportUsageAttribute{}, ReportUsageAttribute{}, "", 2},
	}

	var tests []reportUsageTestCase
	// Create the expected json
	for _, test := range preTests {
		// Check if at least one of the structs isn't empty
		switch {
		case test.serverSize != (ReportUsageAttribute{}) && test.serviceId != (ReportUsageAttribute{}):
			test.expectedResult = fmt.Sprintf(jsonPatterns[test.jsonPatternNum], test.productId, test.commandName, test.serviceId.AttributeName, test.serviceId.AttributeValue, test.serverSize.AttributeName, test.serverSize.AttributeValue)
		case test.serverSize != (ReportUsageAttribute{}):
			test.expectedResult = fmt.Sprintf(jsonPatterns[test.jsonPatternNum], test.productId, test.commandName, test.serverSize.AttributeName, test.serverSize.AttributeValue)
		case test.serviceId != (ReportUsageAttribute{}):
			test.expectedResult = fmt.Sprintf(jsonPatterns[test.jsonPatternNum], test.productId, test.commandName, test.serviceId.AttributeName, test.serviceId.AttributeValue)
		default:
			test.expectedResult = fmt.Sprintf(jsonPatterns[test.jsonPatternNum], test.productId, test.commandName)
		}

		tests = append(tests, test)
	}

	for _, test := range tests {
		t.Run(test.commandName, func(t *testing.T) {
			body, err := reportUsageToJson(test.productId, test.commandName, test.serviceId, test.serverSize)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, string(body))
		})
	}
}
