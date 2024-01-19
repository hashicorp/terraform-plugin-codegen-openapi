## 0.3.0 (January 19, 2024)

ENHANCEMENTS:

* Added data source and resource support for query and path parameters specified in the [OAS Path Item](https://spec.openapis.org/oas/v3.1.0#path-item-object) ([#114](https://github.com/hashicorp/terraform-plugin-codegen-openapi/issues/114))

BUG FIXES:

* Fixed a bug where schema defaults were not detected for integer/int32 properties ([#111](https://github.com/hashicorp/terraform-plugin-codegen-openapi/issues/111))

## 0.2.1 (December 13, 2023)

BUG FIXES:

* Fixed a bug where schemas that used `additionalProperties` with schema composition (allOf/anyOf/oneOf) would return an empty single nested attribute. Will now return map or map nested attribute. ([#100](https://github.com/hashicorp/terraform-plugin-codegen-openapi/issues/100))

## 0.2.0 (October 30, 2023)

FEATURES:

* Added schema.ignores option to generator config for resources, data sources, and providers. Allows excluding attributes from OAS mapping ([#81](https://github.com/hashicorp/terraform-plugin-codegen-openapi/issues/81))

ENHANCEMENTS:

* Added data source support for response body arrays ([#16](https://github.com/hashicorp/terraform-plugin-codegen-openapi/issues/16))
* Schemas that have the `properties` keyword defined with no type will now default to `object` ([#79](https://github.com/hashicorp/terraform-plugin-codegen-openapi/issues/79))

## 0.1.0 (October 17, 2023)

NOTES:

* Initial release of `tfplugingen-openapi` CLI for Terraform Provider Code Generation tech preview ([#68](https://github.com/hashicorp/terraform-plugin-codegen-openapi/issues/68))

