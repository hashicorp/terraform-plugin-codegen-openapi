package explorer

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var _ Explorer = guesstimatorExplorer{}

// This regex is identifies if an API path contains a parameter, indicated by surrounding curly braces
//   - Example: /users/{username} = MATCH
var pathParameterRegex = regexp.MustCompile(`{.*}`)

type guesstimatorExplorer struct {
	spec high.Document
}

type resourceOperations struct {
	// IdentityOps are operations (GET, PUT, POST, DELETE, etc.) on a path that ends with a parameter: /path/{id}
	IdentityOps map[string]*high.Operation

	// CollectionOps are operations (GET, PUT, POST, DELETE, etc.) on a path that don't end with a parameter: /path
	CollectionOps map[string]*high.Operation
}

// As the name suggests, the Guesstimator evaluates an OpenAPIv3 spec and will return
// Resources, DataSources, and their respective names, based on RESTful conventions:
// https://swagger.io/resources/articles/best-practices-in-api-design/
//
// FindResources will group API paths together into collection operations and identity operations, then use the HTTP method to
// determine how to map to a terraform resource. A valid Resource will have a POST collection operation, GET identity operation, and
// a DELETE identity operation. The name of the Resource is a combination of the preceding paths, excluding any path parameters.
// An example of a valid Resource would be:
//   - POST /org/{org_id}/users = Create operation for `org_users` resource
//   - GET /org/{org_id}/users/{id} = Read operation for `org_users` resource
//   - PUT /org/{org_id}/users/{id} = Update operation for `org_users` resource
//   - DELETE /org/{org_id}/users/{id} = Delete operation for `org_users` resource
//
// FindDataSources will group API paths together into collection operations and identity operations, then use the HTTP method to
// determine how to map to a terraform data source. A valid DataSource has a GET identity operation or a GET collection operation.
// The name of the DataSource is a combination of the preceding paths, excluding any path parameters, with an added suffix of "_collection"
// for the collection operation of a DataSource.
// An example of two valid DataSources would be:
//   - GET /org/{org_id}/users = Read operation for `org_users_collection` data source
//   - GET /org/{org_id}/users/{id} = Read operation for `org_users` data source
func NewGuesstimatorExplorer(spec high.Document) Explorer {
	return guesstimatorExplorer{
		spec: spec,
	}
}

func (e guesstimatorExplorer) FindProvider() (Provider, error) {
	// TODO: not sure the best place to automatically pull the provider name... Info section?
	// https://spec.openapis.org/oas/latest.html#info-object
	return Provider{
		Name: "guesstimator_placeholder",
	}, nil
}

// Resource behavior:
// https://developer.hashicorp.com/terraform/language/resources/behavior#how-terraform-applies-a-configuration
func (e guesstimatorExplorer) FindResources() (map[string]Resource, error) {
	resourcesMap := map[string]Resource{}

	groupedResourceOperations := e.groupPathItems()
	for name, group := range groupedResourceOperations {
		if group.CollectionOps == nil || group.IdentityOps == nil {
			continue
		}

		if group.IdentityOps["get"] == nil || group.IdentityOps["delete"] == nil || group.CollectionOps["post"] == nil {
			continue
		}

		// Fallback to POST on identity
		createOp := group.CollectionOps["post"]
		if createOp == nil {
			createOp = group.IdentityOps["post"]
		}

		resourcesMap[name] = Resource{
			CreateOp: createOp,
			ReadOp:   group.IdentityOps["get"],
			UpdateOp: group.IdentityOps["put"],
			DeleteOp: group.IdentityOps["delete"],
		}
	}

	return resourcesMap, nil
}

// Data Source behavior:
// https://developer.hashicorp.com/terraform/language/data-sources#data-resource-behavior
func (e guesstimatorExplorer) FindDataSources() (map[string]DataSource, error) {
	dataSourcesMap := map[string]DataSource{}

	groupedResourceOperations := e.groupPathItems()
	for name, group := range groupedResourceOperations {
		if group.CollectionOps == nil || group.IdentityOps == nil {
			continue
		}

		if group.IdentityOps["get"] != nil {
			// Combine all schemas into something that can be translated to framework IR
			dataSourcesMap[name+"_by_id"] = DataSource{ReadOp: group.IdentityOps["get"]}
		}

		if group.CollectionOps["get"] != nil {
			dataSourcesMap[name+"_collection"] = DataSource{ReadOp: group.CollectionOps["get"]}
		}
	}

	return dataSourcesMap, nil
}

// groupPathItems groups all operations for potential TF resource/data source
//   - Name of resource is determined by combining all nested path segments with underscores
func (e guesstimatorExplorer) groupPathItems() map[string]resourceOperations {
	groups := map[string]resourceOperations{}

	for name, pathItem := range e.spec.Paths.PathItems {
		resource, isIdentity := convertPathToResourceName(name)

		_, ok := groups[resource]
		if !ok {
			groups[resource] = resourceOperations{
				IdentityOps:   map[string]*high.Operation{},
				CollectionOps: map[string]*high.Operation{},
			}
		}

		ops := pathItem.GetOperations()
		for opName, op := range ops {
			if isIdentity {
				groups[resource].IdentityOps[opName] = op
			} else {
				groups[resource].CollectionOps[opName] = op
			}
		}
	}

	return groups
}

// TODO: Consider moving this functionality into a Go type / methods, OASPath.HasIdentityToken(), OASPath.ResourceName(), etc.
// convertPathToResourceName takes a given API path, /example/user/{username},
// and converts it to a valid resource name by combining the paths with underscores, i.e. example_user
func convertPathToResourceName(urlPath string) (string, bool) {
	restOfPath, resource := path.Split(urlPath)
	hasPathParam := pathParameterRegex.Match([]byte(resource))
	if hasPathParam {
		restOfPath, resource = path.Split(path.Dir(restOfPath))
	}

	nameParts := []string{}
	pathParts := strings.FieldsFunc(restOfPath, func(r rune) bool { return r == '/' })

	for i := 0; i < len(pathParts); i++ {
		part := pathParts[i]
		if pathParameterRegex.Match([]byte(part)) {
			continue
		}
		nameParts = append(nameParts, part)
	}

	resourcePrefix := strings.Join(nameParts, "_")
	if resourcePrefix != "" {
		resource = fmt.Sprintf("%s_%s", resourcePrefix, resource)
	}

	return resource, hasPathParam
}
