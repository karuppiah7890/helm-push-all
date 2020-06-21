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
		assert.Contains(t, warnings[0], "chart.metadata.name is required")
		assert.Contains(t, warnings[1], "does not appear to be a gzipped archive")
		assert.Contains(t, warnings[2], "seems to be a YAML file, but expected a gzipped archive")
	}

	assert.Nil(t, err)

	actualChartNames := make([]string, 0, len(chartInfos))

	for _, chartInfo := range chartInfos {
		actualChartNames = append(actualChartNames, chartInfo.Name())
	}

	assert.ElementsMatch(t, expectedChartNames, actualChartNames, "list of charts should match")
}
