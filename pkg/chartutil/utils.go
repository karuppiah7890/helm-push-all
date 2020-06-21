package chartutil

import (
	"io/ioutil"
	"path/filepath"

	helmpush_helm "github.com/chartmuseum/helm-push/pkg/helm"
	"helm.sh/helm/v3/pkg/chart"
	v2chartutil "k8s.io/helm/pkg/chartutil"
)

// Warnings represent warning messages
type Warnings []string

// ChartDependencies are the dependencies / requirements
// of a chart
type ChartDependencies struct {
	V3 []*chart.Dependency
	V2 []*v2chartutil.Dependency
}

// ChartInfo is information about the Chart
type ChartInfo struct {
	chart             *helmpush_helm.Chart
	chartRequirements ChartDependencies
	path              string
	name              string
}

// Name returns the chart name
func (c ChartInfo) Name() string {
	return c.name
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
		requirements, err := getChartDependencies(chart)
		if err != nil {
			return nil, nil, err
		}
		charts = append(charts, ChartInfo{
			chart:             chart,
			chartRequirements: requirements,
			path:              chartPath,
			name:              getChartName(chart),
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

func getChartDependencies(chart *helmpush_helm.Chart) (ChartDependencies, error) {
	v3Chart := chart.V3
	if v3Chart != nil {
		return ChartDependencies{
			V3: v3Chart.Metadata.Dependencies,
		}, nil
	}

	v2Chart := chart.V2
	if v2Chart != nil {
		requirements, err := v2chartutil.LoadRequirements(v2Chart)
		if err != nil && err != v2chartutil.ErrRequirementsNotFound {
			return ChartDependencies{}, err
		}

		var dependencies []*v2chartutil.Dependency
		if requirements != nil {
			dependencies = requirements.Dependencies
		}

		return ChartDependencies{
			V2: dependencies,
		}, nil
	}

	return ChartDependencies{}, nil
}
