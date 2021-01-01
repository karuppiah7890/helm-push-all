package chartutil_test

import (
	"testing"

	"github.com/karuppiah7890/helm-push-all/pkg/chartutil"
	"github.com/stretchr/testify/assert"
)

func TestReadCharts(t *testing.T) {
	expectedChartNames := []string{
		"helm2-chart",
		"helm3-chart",
		"simple-v3-chart",
		"helm2-packaged-chart",
		"helm3-packaged-chart",
	}

	chartInfos, warnings, err := chartutil.ReadCharts("testdata/mixed-charts")
	if assert.Len(t, warnings, 3) {
		assert.Equal(t, "testdata/mixed-charts/invalid-chart: validation: chart.metadata.name is required", warnings[0])
		assert.Equal(t, "testdata/mixed-charts/random-file: file 'testdata/mixed-charts/random-file' does not appear to be a gzipped archive; got 'application/octet-stream'", warnings[1])
		assert.Equal(t, "testdata/mixed-charts/some-yaml.yaml: file 'testdata/mixed-charts/some-yaml.yaml' seems to be a YAML file, but expected a gzipped archive", warnings[2])
	}

	assert.Nil(t, err)

	actualChartNames := make([]string, 0, len(chartInfos))

	for _, chartInfo := range chartInfos {
		actualChartNames = append(actualChartNames, chartInfo.Name())
	}

	assert.ElementsMatch(t, expectedChartNames, actualChartNames, "list of charts should match")
}
