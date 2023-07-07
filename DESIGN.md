<!-- 
    This doc is a placeholder for future documentation pages/markdown files with the intention to surface the assumptions the generator makes (currently in RFC). CLI documentation + generator config syntax should be documented elsewhere
-->

<!--
    TODO: Will need to update this document's wording if the tool changes from generating IR to generating Go code w/ IR library
    TODO: Need to update links to generator config file syntax when available
-->
# OpenAPI to Framework IR (Intermediate Representation) Generator Design

## Overview

The OpenAPI to Framework IR Generator (referred to in this documentation as **generator**) provides mapping between [OpenAPI Specification](https://www.openapis.org/) version 3.0 and 3.1, to Terraform Plugin Framework IR. This mapping currently includes resources, data sources, and provider schema information, all of which are identified with a [generator config file](./README.md).

As the OpenAPI specification (OAS) is designed to describe HTTP APIs in general, it doesn't have full parity with the Terraform Plugin Framework schema or code patterns. There are pieces of logic in the generator that make assumptions on what portions of the OAS to use when mapping to Framework IR, this design document intends to describe those assumptions in detail.

Users of the generator can adjust their OAS to match these assumptions, or suggest changes/customization via the [generator config file](./README.md).

## Determining the OAS Schema to map from operations

### Provider
For generating Provider schema code, the [generator config file](./README.md) defines:
- `provider.name` - required property, which is directly copied to the Framework IR as the name of the Provider.
- `provider.schema_ref` - optional property, which is a [JSON schema reference](https://json-schema.org/understanding-json-schema/structuring.html#ref) to an existing schema in your OpenAPI spec, typically in the [`components.schema` section](https://spec.openapis.org/oas/v3.1.0#fixed-fields-5). This will be used to [map](#mapping-oas-schema-to-plugin-framework-types) the Provider's schema to Framework IR.
```yml
provider:
  name: fakeprovider
  # This schema needs to exist in the OpenAPI spec!
  schema_ref: '#/components/schemas/fake_provider_schema'
```

### Resources
The [generator config file](./README.md) defines the CRUD (`Create`, `Read`, `Update`, `Delete`) operations for a resource in an OAS. In those operations, the generator will search `Create` and `Read` operations for schemas to map to Framework IR. Multiple schemas will be [deep merged](#deep-merge-of-schemas-resources) and the final result will be the Resource schema represented in Framework IR.

```yml
resources:
  fake_thing:
    # Required
    create:
      path: /thing
      method: POST
    read:
      path: /thing/{id}
      method: GET
    # Optional (currently, no effect)
    update:
      path: /thing
      method: PUT
    delete:
      path: /thing/{id}
      method: DELETE
```

#### OAS Schema order (resources)
- `Create` operation [requestBody](https://spec.openapis.org/oas/v3.1.0#requestBodyObject)
    - `requestBody` is the only schema **required** for resources, if not present will log a warning and skip the resource without mapping.
    - Will attempt to use `application/json` first, then will grab the first content-type if not found (alphabetical sort)
- `Create` operation [response](https://spec.openapis.org/oas/v3.1.0#responsesObject)
    - Will attempt to use `200` or `201` first, then will grab the first 2xx response code if not found (lexicographic sort)
    - Will attempt to use `application/json` first, then will grab the first content-type if not found (alphabetical sort)
- `Read` operation [response](https://spec.openapis.org/oas/v3.1.0#responsesObject)
    - Will attempt to use `200` or `201` first, then will grab the first 2xx response code if not found (lexicographic sort)
    - Will attempt to use `application/json` first, then will grab the first content-type if not found (alphabetical sort)
- `Read` operation [parameters](https://spec.openapis.org/oas/v3.1.0#parameterObject)
    - The generator will [deep merge](#deep-merge-of-schemas-resources) the parameters defined at the root of the schema.

#### Deep merge of schemas (resources)
All schemas found will be deep merged together, with the `requestBody` schema from the `Create` operation being the `main schema` that the others will be merged on top. The deep merge has the following characteristics:

- Only attribute name is compared, if the attribute doesn't already exist in the main schema, it will be added. Any mismatched types of the same name will not raise an error and priority will favor the `main schema`.
- Names are strictly compared, so `id` and `user_id` would be two separate attributes in a schema.
- Arrays and Objects will have their child properties merged, so `example_object.string_field` and `example_object.bool_field` will be merged into the same `SingleNestedAttribute` schema.

### Data Sources
The [generator config file](./README.md) defines the `Read` operation for a data source in an OAS. In that operation, the generator will search for a response body schema to map to Framework IR. The response body will be [deep merged](#deep-merge-of-schemas-data-sources) with the query parameters and path parameters of the same `Read` operation and the final result will be the Data Source schema represented in Framework IR.

```yml
data_sources:
  fake_thing:
    read:
      path: /thing/{id}
      method: GET
```

#### OAS Schema order (data sources)
- `Read` operation [response](https://spec.openapis.org/oas/v3.1.0#responsesObject)
    - `response` is the only schema **required** for data sources, if not present will log a warning and skip the data source without mapping.
    - Will attempt to use `200` or `201` first, then will grab the first 2xx response code if not found (lexicographic sort)
    - Will attempt to use `application/json` first, then will grab the first content-type if not found (alphabetical sort)
- `Read` operation [parameters](https://spec.openapis.org/oas/v3.1.0#parameterObject)
    - The generator will [deep merge](#deep-merge-of-schemas-data-sources) the parameters defined at the root of the schema.

#### Deep merge of schemas (data sources)
The response body schema found will be deep merged with the query/path `parameters`, with the `parameters` being the `main schema` that the others will be merged on top. The deep merge has the following characteristics:

- Only attribute name is compared, if the attribute doesn't already exist in the main schema, it will be added. Any mismatched types of the same name will not raise an error and priority will favor the `main schema`.
- Names are strictly compared, so `id` and `user_id` would be two separate attributes in a schema.
- Arrays and Objects will have their child properties merged, so `example_object.string_field` and `example_object.bool_field` will be merged into the same `SingleNestedAttribute` schema.

## Mapping OAS Schema to Plugin Framework Types

### OAS to Plugin Framework Attribute Types

For a given [OAS type](https://spec.openapis.org/oas/v3.1.0#data-types) and format combination, the following rules will be applied for mapping to Framework  attribute types. Not all Framework types are represented natively with OAS, those types are noted below in [Unsupported Attribute Types](#unsupported-attribute-types).

> **NOTE:** All `Type` and `Format` fields below are native to OpenAPI Spec 3.x, with the exception of the format `set`, which is a custom field that only this generator tool is expected to support.

| Type (OAS) | Format (OAS)        | Other Criteria                               | Plugin Framework Attribute Type                                                             |
|------------|---------------------|----------------------------------------------|---------------------------------------------------------------------------------------------|
| `integer`  | -                   | -                                            | `Int64Attribute`                                                                            |
| `number`   | `double` or `float` | -                                            | `Float64Attribute`                                                                          |
| `number`   | -                   | -                                            | `NumberAttribute`                                                                           |
| `string`   | -                   | -                                            | `StringAttribute`                                                                           |
| `boolean`  | -                   | -                                            | `BoolAttribute`                                                                             |
| `array`    | -                   | `items.type == object`                       | `ListNestedAttribute`                                                                       |
| `array`    | -                   | `items.type == (any)`                        | `ListAttribute` (nests with [element types](#oas-to-plugin-framework-element-types))        |
| `array`    | `set`               | `items.type == object`                       | `SetNestedAttribute`                                                                        |
| `array`    | `set`               | `items.type == (any)`                        | `SetAttribute` (nests with [element types](#oas-to-plugin-framework-element-types))         |
| `object`   | -                   | `additionalProperties.type == object`        | `MapNestedAttribute`                                                                        |
| `object`   | -                   | `additionalProperties.type == (any)`         | `MapAttribute`  (nests with [element types](#oas-to-plugin-framework-element-types))        |
| `object`   | -                   | -                                            | `SingleNestedAttribute`                                                                     |

#### Unsupported Attribute Types
- `ListNestedBlock`, `SetNestedBlock`, and `SingleNestedBlock`
    - While the Plugin Framework supports blocks, the Plugin Framework team encourages provider developers to prefer `ListNestedAttribute`, `SetNestedAttribute`, and `SingleNestedAttribute` for new provider development.
- `ObjectAttribute`
    - The generator will default to `SingleNestedAttribute` for object types to provide the additional schema information.

### OAS to Plugin Framework Element Types

For attributes that don't have additional schema information (`ListAttribute`, `SetAttribute`, and `MapAttribute`), the following rules will be applied for mapping from OAS type and format combinations, into Framework element types.

| Type (OAS) | Format (OAS)        | Other Criteria                        | Plugin Framework Element Type   |
|------------|---------------------|---------------------------------------|---------------------------------|
| `integer`  | -                   | -                                     | `Int64Type`                     |
| `number`   | `double` or `float` | -                                     | `Float64Type`                   |
| `number`   | -                   | -                                     | `NumberType`                    |
| `string`   | -                   | -                                     | `StringType`                    |
| `boolean`  | -                   | -                                     | `BoolType`                      |
| `array`    | -                   | -                                     | `ListType`                      |
| `array`    | `set`               | -                                     | `SetType`                       |
| `object`   | -                   | `additionalProperties.type == (any)`  | `MapType`                       |
| `object`   | -                   | -                                     | `ObjectType`                    |

### Required, Computed, and Optional

#### Provider
For the provider, all fields in the provided JSON schema (`provider.schema_ref`) marked as [required](https://json-schema.org/understanding-json-schema/reference/object.html#required-properties) will be mapped as a [Required](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#required) attribute.

If not required, then the field will be mapped as [Optional](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#optional).

#### Resources
For resources, all fields, in the `Create` operation `requestBody` OAS schema, marked as [required](https://json-schema.org/understanding-json-schema/reference/object.html#required-properties) will be mapped as a [Required](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#required) attribute.

If not required, or if the field is in a different schema than the `Create` operation `requestbody`, then the field will be mapped as [Computed](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#computed) and [Optional](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#optional).

#### Data Sources
For data sources, all fields, in the `Read` operation `parameters` OAS schema, marked as [required](https://json-schema.org/understanding-json-schema/reference/object.html#required-properties) will be mapped as a [Required](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#required) attribute.

If not required, or if the field is in a different schema than the `Read` operation `parameters`, then the field will be mapped as [Computed](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#computed) and [Optional](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#optional).

### Other field mapping

| Field (OAS)                                                                   | Field (Plugin Framework Schema)                                                                                                           |
|-------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| [default](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-default) | [(Attribute).Default](https://developer.hashicorp.com/terraform/plugin/framework/resources/default) (resources only) |
| [description](https://spec.openapis.org/oas/latest.html#rich-text-formatting) | [(Attribute).MarkdownDescription](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#markdowndescription-1) |
| [enum](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-enum) | [(Attribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |
| [format (password)](https://spec.openapis.org/oas/latest.html#data-types)     | [(StringAttribute).Sensitive](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#sensitive)                 |
| [maximum](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-maximum) | [(Attribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |
| [maxItems](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-maxItems) | [(Attribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |
| [maxLength](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-maxLength) | [(StringAttribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |
| [maxProperties](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-maxProperties) | [(MapAttribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |
| [minimum](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-minimum) | [(Attribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |
| [minItems](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-minItems) | [(Attribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |
| [minLength](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-minLength) | [(StringAttribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |
| [minProperties](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-minProperties) | [(Attribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |
| [pattern](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-pattern) | [(StringAttribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |
| [uniqueItems](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-uniqueItems) | [(ListAttribute).Validators](https://developer.hashicorp.com/terraform/plugin/framework/validation) |

## Multi-type Support

Generally, [multi-types](https://cswr.github.io/JsonSchema/spec/multiple_types/) are not supported by the generator as the Terraform Plugin Framework does not support multi-types. There is one specific scenario that is supported by the generator and that is any type that is combined with the `null` type, as any Plugin Framework attribute can hold a [null](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes#null) type.

### Nullable Multi-type support
> **Note:** with nullable multi-types, the `description` will be populated from the root-level schema, as shown below. 

In an OAS schema, the following keywords defining nullable multi-types are supported (nullable types will follow the same mapping rules [defined above](#oas-to-plugin-framework-attribute-types) for the type that is not the `null` type):

#### `type` keyword array
```jsonc
// Maps to StringAttribute
{
  "nullable_string_example": {
    "description": "this is the description that's used!",
    "type": [
      "string",
      "null"
    ]
  }
}

// Maps to Int64Attribute
{
  "nullable_integer_example": {
    "description": "this is the description that's used!",
    "type": [
      "null",
      "integer"
    ]
  }
}
```

#### `anyOf` and `oneOf` keywords
```jsonc
// Maps to SingleNestedAttribute
{
  "nullable_object_one": {
    "description": "this is the description that's used!",
    "anyOf": [
      {
        "type": "null"
      },
      {
        "$ref": "#/components/schemas/example_object_one"
      }
    ]
  }
}

// Maps to SingleNestedAttribute
{
  "nullable_object_two": {
    "description": "this is the description that's used!",
    "oneOf": [
      {
        "$ref": "#/components/schemas/example_object_two"
      },
      {
        "type": "null"
      }
    ]
  }
}
```
