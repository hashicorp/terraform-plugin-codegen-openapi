package explorer

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/config"
	"strings"

	high "github.com/pb33f/libopenapi/datamodel/high/v3"
	low "github.com/pb33f/libopenapi/datamodel/low/v3"
)

var _ Explorer = configExplorer{}

type configExplorer struct {
	spec   high.Document
	config config.Config
}

// A ConfigExplorer will use an additional config file to identify resource and data source operations in a provided
// OpenAPIv3 spec. This additional config file will provide information such as:
//   - Create/Read/Update/Delete endpoints/URLs (schema will be automatically grabbed via request/response body and parameters in mapper)
//   - Resource + Data Source names
func NewConfigExplorer(spec high.Document, cfg config.Config) Explorer {
	return configExplorer{
		spec:   spec,
		config: cfg,
	}
}

func (e configExplorer) FindProvider() (Provider, error) {
	return Provider{
		Name: e.config.Provider.Name,
	}, nil
}

func (e configExplorer) FindResources() (map[string]Resource, error) {
	resources := map[string]Resource{}
	for name, opMapping := range e.config.Resources {
		// TODO: should we throw an error if an invalid or non-existent path/methods are given?
		resources[name] = Resource{
			CreateOp: extractOp(e.spec.Paths, opMapping.Create),
			ReadOp:   extractOp(e.spec.Paths, opMapping.Read),
			UpdateOp: extractOp(e.spec.Paths, opMapping.Update),
			DeleteOp: extractOp(e.spec.Paths, opMapping.Delete),
		}
	}

	return resources, nil
}

func (e configExplorer) FindDataSources() (map[string]DataSource, error) {
	dataSources := map[string]DataSource{}
	for name, opMapping := range e.config.DataSources {
		dataSources[name] = DataSource{
			ReadOp: extractOp(e.spec.Paths, opMapping.Read),
		}
	}
	return dataSources, nil
}

func extractOp(paths *high.Paths, oasLocation *config.OpenApiSpecLocation) *high.Operation {
	if oasLocation == nil || paths == nil || paths.PathItems == nil || paths.PathItems[oasLocation.Path] == nil {
		return nil
	}

	switch strings.ToLower(oasLocation.Method) {
	case low.PostLabel:
		return paths.PathItems[oasLocation.Path].Post
	case low.GetLabel:
		return paths.PathItems[oasLocation.Path].Get
	case low.PutLabel:
		return paths.PathItems[oasLocation.Path].Put
	case low.DeleteLabel:
		return paths.PathItems[oasLocation.Path].Delete
	case low.PatchLabel:
		return paths.PathItems[oasLocation.Path].Patch
	case low.OptionsLabel:
		return paths.PathItems[oasLocation.Path].Options
	case low.HeadLabel:
		return paths.PathItems[oasLocation.Path].Head
	case low.TraceLabel:
		return paths.PathItems[oasLocation.Path].Trace
	default:
		return nil
	}
}
