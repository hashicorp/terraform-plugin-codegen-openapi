package explorer

import high "github.com/pb33f/libopenapi/datamodel/high/v3"

var _ Explorer = configExplorer{}

type configExplorer struct {
	spec high.Document
	// config TBD
}

func (e configExplorer) FindResources() (map[string]Resource, error) {
	panic("unimplemented")
}

func (e configExplorer) FindDataSources() (map[string]DataSource, error) {
	panic("unimplemented")
}

// A ConfigExplorer will use an additional config file to identify resource and data source operations in a provided
// OpenAPIv3 spec. This additional config file will provide information such as:
//   - Create/Read/Update/Delete endpoints/URLs (schema will still be automatically grabbed via request body/response in mapper)
//   - Resource + Data Source names
//   - Additional customization not yet defined
//
// This will most likely be the default mode of operation, as the Guesstimator Explorer in it's current state is not accurate
func NewConfigExplorer(spec high.Document) Explorer {
	return configExplorer{
		spec: spec,
	}
}
