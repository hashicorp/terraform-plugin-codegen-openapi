# OpenAPI Provider Spec Generator Design

## Mapping OAS to Provider Code Specification

The rules for mapping an OpenAPI spec (OAS) to [Provider Code Specification](https://developer.hashicorp.com/terraform/plugin/code-generation/specification).

### Provider
For generating the [Provider](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#provider) specification, the generator config defines a single `provider` object:
```yml
provider:
  name: examplecloud
  # This schema needs to exist in the OpenAPI spec!
  schema_ref: '#/components/schemas/examplecloud_provider_schema'
```

- `name` is directly copied to the provider code specification field: `provider.name`.
- `schema_ref` is a [JSON schema reference](https://json-schema.org/understanding-json-schema/structuring.html#ref) that is used to [map](#oas-types-to-provider-attributes) to the Provider's schema: `provider.schema`


### Resources

For generating [Resource](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#resource) specifications, the generator config defines a map `resources`:

```yml
resources:
  thing:
    create:
      path: /thing
      method: POST
    read:
      path: /thing/{id}
      method: GET
    update:
      path: /thing
      method: PUT
    delete:
      path: /thing/{id}
      method: DELETE
```

In these OAS operations, the generator will search the `create` and `read` for schemas to map to the provider code specification. Multiple schemas will have the [OAS types mapped to Provider Attributes](#oas-types-to-provider-attributes) and then be merged together; with the final result being the [Resource](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#resource) `schema`. The schemas that will be merged together (in priority order):
1. `create` operation: [requestBody](https://spec.openapis.org/oas/v3.1.0#requestBodyObject)
    - `requestBody` is the only schema **required** for resources. If not found, the generator will skip the resource without mapping.
    - Will attempt to use `application/json` content-type first. If not found, will grab the first available content-type with a schema (alphabetical order)
2. `create` operation: response body in [responses](https://spec.openapis.org/oas/v3.1.0#responsesObject)
    - Will attempt to use `200` or `201` response body. If not found, will grab the first available `2xx` response code with a schema (lexicographic order)
    - Will attempt to use `application/json` content-type first. If not found, will grab the first available content-type with a schema (alphabetical order)
3. `read` operation: response body in [responses](https://spec.openapis.org/oas/v3.1.0#responsesObject)
    - Will attempt to use `200` or `201` response body. If not found, will grab the first available `2xx` response code with a schema (lexicographic order)
    - Will attempt to use `application/json` content-type first. If not found, will grab the first available content-type with a schema (alphabetical order)
4. `read` operation: [parameters](https://spec.openapis.org/oas/v3.1.0#parameterObject)
    - The generator will merge all `query` and `path` parameters to the root of the schema.

All schemas found will be deep merged together, with the `requestBody` schema from the `create` operation being the **main schema** that the others will be merged on top. The deep merge has the following characteristics:

- Only attribute name is compared, if the attribute doesn't already exist in the **main schema**, it will be added. Any mismatched types of the same name will not raise an error and priority will favor the **main schema**.
	- Names are strictly compared, so `id` and `user_id` would be two separate attributes in a schema.
- Arrays and Objects will have their child attributes merged, so `example_object.string_field` and `example_object.bool_field` will be merged into the same `SingleNestedAttribute` schema.

### Data Sources

For generating [Data Source](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#data-source) specifications, the generator config defines a map `data_sources`:

```yml
data_sources:
  thing:
    read:
      path: /thing/{id}
      method: GET
```

The generator uses the `read` operation to map to the provider code specification. Multiple schemas will have the [OAS types mapped to Provider Attributes](#oas-types-to-provider-attributes) and then be merged together; with the final result being the [Data Source](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#data-source) `schema`. The schemas that will be merged together (in priority order):
1. `read` operation: [parameters](https://spec.openapis.org/oas/v3.1.0#parameterObject)
    - The generator will merge all `query` and `path` parameters to the root of the schema.
2. `read` operation: response body in [responses](https://spec.openapis.org/oas/v3.1.0#responsesObject)
    - The response body is the only schema **required** for data sources. If not found, the generator will skip the data source without mapping.
    - Will attempt to use `200` or `201` response body. If not found, will grab the first available `2xx` response code with a schema (lexicographic order)
    - Will attempt to use `application/json` content-type first. If not found, will grab the first available content-type with a schema (alphabetical order)

The response body schema found will be deep merged with the query/path `parameters`, with the `parameters` being the **main schema** that the others will be merged on top. The deep merge has the following characteristics:

- Only attribute name is compared, if the attribute doesn't already exist in the **main schema**, it will be added. Any mismatched types of the same name will not raise an error and priority will favor the **main schema**.
  - Names are strictly compared, so `id` and `user_id` would be two separate attributes in a schema.
- Arrays and Objects will have their child attributes merged, so `example_object.string_field` and `example_object.bool_field` will be merged into the same `SingleNestedAttribute` schema.

#### Collection Data Sources

If the response body schema for a data source is of type `array`, the schema in `items` will be mapped to a collection attribute (`ListNested`, `SetNested`, `List`, `Set`) at the root of the mapped data source. The name of the attribute will be the same as the data source name from the generator config. All [mapping rules](#oas-types-to-provider-attributes) will be followed for nested attributes.

##### Generator Config
```yaml
provider:
  name: petstore

data_sources:
  pets:
    read:
      path: /pet/findByStatus
      method: GET
```

##### OpenAPI Spec
```jsonc
{
  // ... Rest of OAS
  "/pet/findByStatus": {
    "get": {
      "responses": {
        "200": {
          "description": "successful operation",
          "content": {
            "application/json": {
              "schema": {
                "type": "array",
                "items": {
                  "$ref": "#/components/schemas/Pet"
                }
              }
            }
          }
        }
      }
    }
  }
}
```

##### Provider Code Spec output
```jsonc
{
  "datasources": [
    {
      "name": "pets",
      "schema": {
        "attributes": [
          {
            "name": "pets",
            "list_nested": {
              "computed_optional_required": "computed",
              "nested_object": {
                "attributes": [
                  // ... mapping of #/components/schemas/Pet
                ]
              }
            }
          }
        ]
      }
    }
  ],
  "provider": {
    "name": "petstore"
  }
}
```



### OAS Types to Provider Attributes

For a given OAS [`type`](https://spec.openapis.org/oas/v3.1.0#data-types) and `format` combination, the following rules will be applied for mapping to the provider code specification. Not all Provider attributes are represented natively with OAS, those types are noted below in [Unsupported Attributes](#unsupported-attributes).

<Note>
All <b>Type</b> and <b>Format</b> fields below are native to OpenAPI Spec 3.x, with the exception of the format <b>set</b>, which is a custom format that only this generator tool is expected to support.
</Note>

| Type (OAS) | Format (OAS)        | Other Criteria                               | Provider Attribute Type                                                                     |
|------------|---------------------|----------------------------------------------|---------------------------------------------------------------------------------------------|
| `boolean`  | -                   | -                                            | `BoolAttribute`                                                                             |
| `integer`  | -                   | -                                            | `Int64Attribute`                                                                            |
| `number`   | `double` or `float` | -                                            | `Float64Attribute`                                                                          |
| `number`   | -                   | -                                            | `NumberAttribute`                                                                           |
| `string`   | -                   | -                                            | `StringAttribute`                                                                           |
| `array`    | -                   | `items.type == object`                       | `ListNestedAttribute`                                                                       |
| `array`    | -                   | `items.type == (any)`                        | `ListAttribute` (nests with [element types](#oas-types-to-provider-element-types))          |
| `array`    | `set`               | `items.type == object`                       | `SetNestedAttribute`                                                                        |
| `array`    | `set`               | `items.type == (any)`                        | `SetAttribute` (nests with [element types](#oas-types-to-provider-element-types))           |
| `object`   | -                   | `additionalProperties.type == object`        | `MapNestedAttribute`                                                                        |
| `object`   | -                   | `additionalProperties.type == (any)`         | `MapAttribute`  (nests with [element types](#oas-types-to-provider-element-types))          |
| `object`   | -                   | -                                            | `SingleNestedAttribute`                                                                     |

#### Unsupported Attributes
- `ListNestedBlock`, `SetNestedBlock`, and `SingleNestedBlock`
    - While the provider code specification supports blocks, the recommendation is to prefer `ListNestedAttribute`, `SetNestedAttribute`, and `SingleNestedAttribute` for new provider development.
- `ObjectAttribute`
    - The generator will default to `SingleNestedAttribute` for object types to provide additional schema information.

#### OAS Types to Provider Element Types

For attributes that don't have additional schema information (`ListAttribute`, `SetAttribute`, and `MapAttribute`), the following rules will be applied for mapping from an OAS `type` and `format` combination, into Provider element types.

| Type (OAS) | Format (OAS)        | Other Criteria                        | Provider Element Type           |
|------------|---------------------|---------------------------------------|---------------------------------|
| `boolean`  | -                   | -                                     | `BoolType`                      |
| `integer`  | -                   | -                                     | `Int64Type`                     |
| `number`   | `double` or `float` | -                                     | `Float64Type`                   |
| `number`   | -                   | -                                     | `NumberType`                    |
| `string`   | -                   | -                                     | `StringType`                    |
| `array`    | -                   | -                                     | `ListType`                      |
| `array`    | `set`               | -                                     | `SetType`                       |
| `object`   | -                   | `additionalProperties.type == (any)`  | `MapType`                       |
| `object`   | -                   | -                                     | `ObjectType`                    |

#### Provider - Required or Optional
For the provider, all fields in the provided JSON schema (`provider.schema_ref`) marked as [required](https://json-schema.org/understanding-json-schema/reference/object.html#required-properties) will be mapped as `required`.

If not required, then the field will be mapped as `optional`.

#### Resources - Required, Computed or Optional
For resources, all fields in the `create` operation `requestBody` OAS schema marked as [required](https://json-schema.org/understanding-json-schema/reference/object.html#required-properties) will be mapped as `required`. If [default](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-default) is also specified, it will be mapped as `computed_optional` instead.

If not required, then the field will be mapped as `computed_optional`.

If the field is only present in a schema other than the `create` operation `requestBody`, then the field will be mapped as `computed`.

#### Data Sources - Required, Computed or Optional
For data sources, all fields in the `read` operation `parameters` OAS schema marked as [required](https://json-schema.org/understanding-json-schema/reference/object.html#required-properties) will be mapped as `required`.

If not required, then the field will be mapped as `computed_optional`.

If the field is only present in a schema other than the `read` operation `parameters`, then the field will be mapped as `computed`.

#### Other OAS field mappings

| Field (OAS)                                                                                           | Field ([Provider Code Specification](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#attribute-type)) |
|-------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------|
| [default](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-default)             | [`default`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#default) (resources only)                 |
| [deprecated](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-deprecated)       | `deprecation_message`                                                                                 |
| [description](https://spec.openapis.org/oas/latest.html#rich-text-formatting)                         | `description`                                                                                         |
| [enum](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-enum)                   | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |
| [format (password)](https://spec.openapis.org/oas/latest.html#data-types)                             | `sensitive`                                                                                           |
| [maximum](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-maximum)             | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |
| [maxItems](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-maxItems)           | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |
| [maxLength](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-maxLength)         | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |
| [maxProperties](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-maxProperties) | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |
| [minimum](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-minimum)             | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |
| [minItems](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-minItems)           | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |
| [minLength](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-minLength)         | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |
| [minProperties](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-minProperties) | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |
| [pattern](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-pattern)             | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |
| [uniqueItems](https://json-schema.org/draft/2020-12/json-schema-validation.html#name-uniqueItems)     | [`validators`](https://developer.hashicorp.com/terraform/plugin/code-generation/specification#validators)                            |

### Attribute Names
After all attributes have been [mapped](#oas-types-to-provider-attributes) and any overrides/aliases have been applied, the attribute names mapped from the OAS will be converted (if needed) to valid [Terraform Identifiers](https://developer.hashicorp.com/terraform/language/syntax/configuration#identifiers). This [logic](https://github.com/hashicorp/terraform-plugin-codegen-openapi/blob/main/internal/mapper/util/framework_identifier.go#L25) performs the following, in order:
1. Removes all characters that are NOT alphanumeric or an underscore
2. Removes all leading numbers
3. Inserts an underscore between any lowercase letter that is immediately followed by an uppercase letter
4. Lowercases the final result

See the [test cases](https://github.com/hashicorp/terraform-plugin-codegen-openapi/blob/main/internal/mapper/util/framework_identifier_test.go#L15) for examples on the expectations of this conversion process.

This ensures all properties from an OAS are converted to valid Terraform identifiers, but can technically cause conflicts if multiple distinct OAS properties are scrubbed to the same value:
- `Fake_Thing` -> `fake_thing`
- `fakeThing` -> `fake_thing`

## Known Limitations
As OpenAPI is designed to describe HTTP APIs in general, it doesn't always fully align with [Terraform Provider design principles](https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles). There are pieces of logic in this generator that make assumptions on what portions of the OAS to use when mapping to the provider code specification, however there are some limitations on what can be supported, which are documented below.

### Multi-type Support

Generally, [multi-types](https://cswr.github.io/JsonSchema/spec/multiple_types/) are not supported by the generator as the Terraform Plugin Framework does not support multi-types. There are two specific scenarios that are supported by the generator. 

> **Note:** with multi-type support described below, the `description` will be populated from the root-level schema, see examples. 

### Nullable Multi-type support

If a multi-type is detected where one of the types is `null`, the other type will be used for schema mapping using the same rules [defined above](#oas-types-to-provider-attributes).

#### Examples with `type` array
```json
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

#### Examples with `anyOf` and `oneOf`
```json
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

### String-able Multi-type support

If a multi-type is detected where one of the types is a `string` and the other type is a `primitive`, then the resulting attribute will be a `StringAttribute`.

Supported `primitive` types that can be represented as `string`:
- `number`
- `integer`
- `boolean`


#### Examples with `type` array, `oneOf`, and `anyOf`
```json
// Maps to StringAttribute
{
  "stringable_number_example": {
    "description": "this is the description that's used!",
    "type": [
      "string",
      "number"
    ]
  }
}

// Maps to StringAttribute
{
  "stringable_integer_example": {
    "description": "this is the description that's used!",
    "anyOf": [
      {
        "type": "integer"
      },
      {
        "type": "string"
      }
    ]
  }
}

// Maps to StringAttribute
{
  "stringable_boolean_example": {
    "description": "this is the description that's used!",
    "oneOf": [
      {
        "type": "string"
      },
      {
        "type": "boolean"
      }
    ]
  }
}
```
