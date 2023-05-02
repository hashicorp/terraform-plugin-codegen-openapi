package explorer_test

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/explorer"
	"testing"

	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func Test_GuesstimatorExplorer_FindResources(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		pathItems         map[string]*high.PathItem
		expectedResources []string
	}{
		"valid flat resource combo": {
			pathItems: map[string]*high.PathItem{
				"/resources": {
					Post: &high.Operation{},
				},
				"/resources/{resource_id}": {
					Get:    &high.Operation{},
					Delete: &high.Operation{},
				},
			},
			expectedResources: []string{"resources"},
		},
		"valid nested resource combo": {
			pathItems: map[string]*high.PathItem{
				"/verycool/verynice/resources": {
					Post: &high.Operation{},
				},
				"/verycool/verynice/resources/{resource_id}": {
					Get:    &high.Operation{},
					Delete: &high.Operation{},
				},
			},
			expectedResources: []string{"verycool_verynice_resources"},
		},
		"valid nested with id resource combo": {
			pathItems: map[string]*high.PathItem{
				"/verycool/{id}/verynice/resources": {
					Post: &high.Operation{},
				},
				"/verycool/{id}/verynice/resources/{resource_id}": {
					Get:    &high.Operation{},
					Delete: &high.Operation{},
				},
			},
			expectedResources: []string{"verycool_verynice_resources"},
		},
		"invalid resource combo - POST,DELETEbyID": {
			pathItems: map[string]*high.PathItem{
				"/resources": {
					Post: &high.Operation{},
				},
				"/resources/{resource_id}": {
					Delete: &high.Operation{},
				},
			},
			expectedResources: []string{},
		},
		"invalid resource combo - GETbyID,DELETEbyID": {
			pathItems: map[string]*high.PathItem{
				"/resources/{resource_id}": {
					Get:    &high.Operation{},
					Delete: &high.Operation{},
				},
			},
			expectedResources: []string{},
		},
		"invalid resource combo - GETbyID,POST": {
			pathItems: map[string]*high.PathItem{
				"/resources": {
					Post: &high.Operation{},
				},
				"/resources/{resource_id}": {
					Get: &high.Operation{},
				},
			},
			expectedResources: []string{},
		},
		"invalid resource combo - no ops": {
			pathItems: map[string]*high.PathItem{
				"/resources":               {},
				"/resources/{resource_id}": {},
			},
			expectedResources: []string{},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			explorer := explorer.NewGuesstimatorExplorer(high.Document{Paths: &high.Paths{PathItems: testCase.pathItems}})
			resources, err := explorer.FindResources()

			if err != nil {
				t.Fatalf("was not expecting error, got: %s", err)
			}

			if len(resources) != len(testCase.expectedResources) {
				t.Fatalf("expected %d resources, found %d resources", len(testCase.expectedResources), len(resources))
			}

			for _, expectedResource := range testCase.expectedResources {
				_, ok := resources[expectedResource]
				if !ok {
					t.Fatalf("%s resource not found", expectedResource)
				}
			}
		})
	}
}

func Test_GuesstimatorExplorer_FindDataSources(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		pathItems           map[string]*high.PathItem
		expectedDataSources []string
	}{
		"valid flat data source combo": {
			pathItems: map[string]*high.PathItem{
				"/resources": {
					Get: &high.Operation{},
				},
				"/resources/{resource_id}": {
					Get: &high.Operation{},
				},
			},
			expectedDataSources: []string{"resources_collection", "resources_by_id"},
		},
		"valid nested data source combo": {
			pathItems: map[string]*high.PathItem{
				"/verycool/verynice/resources": {
					Post: &high.Operation{},
				},
				"/verycool/verynice/resources/{resource_id}": {
					Get:    &high.Operation{},
					Delete: &high.Operation{},
				},
			},
			expectedDataSources: []string{"verycool_verynice_resources_by_id"},
		},
		"valid nested with id data source combo": {
			pathItems: map[string]*high.PathItem{
				"/verycool/{id}/verynice/resources": {
					Get:  &high.Operation{},
					Post: &high.Operation{},
				},
				"/verycool/{id}/verynice/resources/{resource_id}": {
					Delete: &high.Operation{},
				},
			},
			expectedDataSources: []string{"verycool_verynice_resources_collection"},
		},
		"invalid data source combo - no matching ops": {
			pathItems: map[string]*high.PathItem{
				"/resources": {
					Put:     &high.Operation{},
					Post:    &high.Operation{},
					Delete:  &high.Operation{},
					Options: &high.Operation{},
					Head:    &high.Operation{},
					Patch:   &high.Operation{},
					Trace:   &high.Operation{},
				},
				"/resources/{resource_id}": {
					Put:     &high.Operation{},
					Post:    &high.Operation{},
					Delete:  &high.Operation{},
					Options: &high.Operation{},
					Head:    &high.Operation{},
					Patch:   &high.Operation{},
					Trace:   &high.Operation{},
				},
			},
			expectedDataSources: []string{},
		},
		"invalid data source combo - no ops": {
			pathItems: map[string]*high.PathItem{
				"/resources":               {},
				"/resources/{resource_id}": {},
			},
			expectedDataSources: []string{},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			explorer := explorer.NewGuesstimatorExplorer(high.Document{Paths: &high.Paths{PathItems: testCase.pathItems}})
			dataSources, err := explorer.FindDataSources()

			if err != nil {
				t.Fatalf("was not expecting error, got: %s", err)
			}

			if len(dataSources) != len(testCase.expectedDataSources) {
				t.Fatalf("expected %d data sources, found %d data sources", len(testCase.expectedDataSources), len(dataSources))
			}

			for _, expectedDataSource := range testCase.expectedDataSources {
				_, ok := dataSources[expectedDataSource]
				if !ok {
					t.Fatalf("%s data sources not found", expectedDataSource)
				}
			}
		})
	}
}
