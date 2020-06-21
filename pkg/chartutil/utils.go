package chartutil

import (
	"io/ioutil"
	"path/filepath"

	helmpush_helm "github.com/chartmuseum/helm-push/pkg/helm"
)

// Warnings represent warning messages
type Warnings []string

// ReadCharts from a directory
func ReadCharts(dir string) ([]*helmpush_helm.Chart, Warnings, error) {
	warnings := Warnings{}
	charts := []*helmpush_helm.Chart{}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	for _, fileInfo := range fileInfos {
		filePath := filepath.Join(dir, fileInfo.Name())
		chart, warning := helmpush_helm.GetChartByName(filePath)
		if warning != nil {
			warnings = append(warnings, warning.Error())
			continue
		}
		charts = append(charts, chart)
	}

	return charts, warnings, nil
}
