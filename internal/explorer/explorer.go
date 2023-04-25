package explorer

import high "github.com/pb33f/libopenapi/datamodel/high/v3"

// Implementations of the Explorer interface will relate OpenAPIv3 operations to a set of Terraform Provider actions
//   - https://spec.openapis.org/oas/latest.html#operation-object
type Explorer interface {
	// TODO: Add Provider in here?
	//
	// Maps all the way down!
	//   _____     ____
	//  /      \  |  o |
	// |  map   |/ ___\|
	// |_________/
	// |_|_| |_|_|
	//   _____     ____
	//  /      \  |  o |
	// |  map   |/ ___\|
	// |_________/
	// |_|_| |_|_|
	//   _____     ____
	//  /      \  |  o |
	// |  map   |/ ___\|
	// |_________/
	// |_|_| |_|_|
	//
	FindResources() (map[string]Resource, error)
	FindDataSources() (map[string]DataSource, error)
}

type Resource struct {
	CreateOp *high.Operation
	ReadOp   *high.Operation
	UpdateOp *high.Operation
	DeleteOp *high.Operation
}

type DataSource struct {
	ReadOp *high.Operation
}
