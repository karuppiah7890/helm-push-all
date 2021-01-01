package chartutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
)

func TestFindRelationship(t *testing.T) {
	abcdChartRepoURL := "https://abcd-chart-repo.com"
	chartA := &ChartInfo{
		chartRequirements: ChartDependencies{
			V3: []*chart.Dependency{
				{
					Name:       "chartB",
					Repository: abcdChartRepoURL,
				},
				{
					Name:       "chartC",
					Repository: abcdChartRepoURL,
				},
			},
		},
	}
	chartB := &ChartInfo{
		chartRequirements: ChartDependencies{
			V3: []*chart.Dependency{
				{
					Name:       "chartD",
					Repository: "https://some-chart-repo.com",
				},
				{
					Name:       "chartE",
					Repository: abcdChartRepoURL,
				},
			},
		},
	}
	chartC := &ChartInfo{
		chartRequirements: ChartDependencies{
			V3: []*chart.Dependency{
				{
					Name:       "chartD",
					Repository: abcdChartRepoURL,
				},
			},
		},
	}
	chartD := &ChartInfo{}
	chartE := &ChartInfo{}

	chartInfos := map[string]*ChartInfo{
		"chartA": chartA,
		"chartB": chartB,
		"chartC": chartC,
		"chartD": chartD,
		"chartE": chartE,
	}

	expectedChartAChildren := []*ChartInfo{chartB, chartC}
	expectedChartBChildren := []*ChartInfo{chartE}
	expectedChartCChildren := []*ChartInfo{chartD}

	expectedChartBParents := []*ChartInfo{chartA}
	expectedChartCParents := []*ChartInfo{chartA}
	expectedChartDParents := []*ChartInfo{chartC}
	expectedChartEParents := []*ChartInfo{chartB}

	FindRelationship(chartInfos, abcdChartRepoURL)

	assert.ElementsMatch(t, expectedChartAChildren, chartA.Children())
	assert.ElementsMatch(t, expectedChartBChildren, chartB.Children())
	assert.ElementsMatch(t, expectedChartCChildren, chartC.Children())
	assert.ElementsMatch(t, nil, chartD.Children())
	assert.ElementsMatch(t, nil, chartE.Children())

	assert.ElementsMatch(t, nil, chartA.Parents())
	assert.ElementsMatch(t, expectedChartBParents, chartB.Parents())
	assert.ElementsMatch(t, expectedChartCParents, chartC.Parents())
	assert.ElementsMatch(t, expectedChartDParents, chartD.Parents())
	assert.ElementsMatch(t, expectedChartEParents, chartE.Parents())
}
