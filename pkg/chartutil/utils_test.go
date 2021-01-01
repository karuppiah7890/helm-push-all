package chartutil_test

import (
	"fmt"
	"path/filepath"
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

	invalidChartPath1 := filepath.Join("testdata", "mixed-charts", "invalid-chart")
	invalidChartPath2 := filepath.Join("testdata", "mixed-charts", "random-file")
	invalidChartPath3 := filepath.Join("testdata", "mixed-charts", "some-yaml.yaml")

	expectedWarnings := chartutil.Warnings{
		fmt.Sprintf("%s: validation: chart.metadata.name is required", invalidChartPath1),
		fmt.Sprintf("%s: file '%s' does not appear to be a gzipped archive; got 'application/octet-stream'", invalidChartPath2, invalidChartPath2),
		fmt.Sprintf("%s: file '%s' seems to be a YAML file, but expected a gzipped archive", invalidChartPath3, invalidChartPath3),
	}

	chartInfos, warnings, err := chartutil.ReadCharts("testdata/mixed-charts")
	if assert.Len(t, warnings, 3) {
		assert.Equal(t, expectedWarnings, warnings)
	}

	assert.Nil(t, err)

	actualChartNames := make([]string, 0, len(chartInfos))

	for _, chartInfo := range chartInfos {
		actualChartNames = append(actualChartNames, chartInfo.Name())
	}

	assert.ElementsMatch(t, expectedChartNames, actualChartNames, "list of charts should match")
}
