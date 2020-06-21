package chartutil

import (
	"io/ioutil"
	"path/filepath"

	helmpush_helm "github.com/chartmuseum/helm-push/pkg/helm"
)

// Warnings represent warning messages
type Warnings []string

// ChartInfo is information about the Chart
type ChartInfo struct {
	Chart *helmpush_helm.Chart
	Path  string
	Name  string
}

// ChartInfos is an array of chart
type ChartInfos []ChartInfo

// ReadCharts from a directory
func ReadCharts(dir string) (ChartInfos, Warnings, error) {
	warnings := Warnings{}
	charts := ChartInfos{}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	for _, fileInfo := range fileInfos {
		chartPath := filepath.Join(dir, fileInfo.Name())
		chart, warning := helmpush_helm.GetChartByName(chartPath)
		if warning != nil {
			warnings = append(warnings, warning.Error())
			continue
		}
		charts = append(charts, ChartInfo{
			Chart: chart,
			Path:  chartPath,
			Name:  getChartName(chart),
		})
	}

	return charts, warnings, nil
}

func getChartName(chart *helmpush_helm.Chart) string {
	v3Chart := chart.V3
	if v3Chart != nil {
		return v3Chart.Name()
	}

	v2Chart := chart.V2
	if v2Chart != nil {
		return v2Chart.Metadata.Name
	}

	return ""
}
