package explorer

import high "github.com/pb33f/libopenapi/datamodel/high/v3"

func mergeParameters(commonParameters []*high.Parameter, operation *high.Operation) []*high.Parameter {
	mergedParameters := commonParameters
	if operation != nil {
		for _, operationParameter := range operation.Parameters {
			found := false
			for i, mergedParameter := range mergedParameters {
				if operationParameter.Name == mergedParameter.Name {
					found = true
					mergedParameters[i] = operationParameter
					break
				}
			}
			if !found {
				mergedParameters = append(mergedParameters, operationParameter)
			}
		}
	}
	return mergedParameters
}

func (e *Resource) ReadOpParameters() []*high.Parameter {
	return mergeParameters(e.Parameters, e.ReadOp)
}

func (e *DataSource) ReadOpParameters() []*high.Parameter {
	return mergeParameters(e.Parameters, e.ReadOp)
}
