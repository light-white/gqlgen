---
linkTitle: Configuration
title: How to configure gqlgen using gqlgen.yml
description: How to configure gqlgen using gqlgen.yml
menu: main
weight: -5
---

gqlgen can be configured using a `gqlgen.yml` file, by default it will be loaded from the current directory, or any parent directory.

As an example, here is the default configuration file generated with `gqlgen init`:

```yml
# Where are all the schema files located? globs are supported eg  src/**/*.graphqls
schema:
  - graph/*.graphqls

# Where should the generated server code go?
exec:
  package: graph
  layout: single-file # Only other option is "follow-schema," ie multi-file.

  # Only for single-file layout:
  filename: graph/generated.go

  # Only for follow-schema layout:
  # dir: graph
  # filename_template: "{name}.generated.go"

  # Optional: Maximum number of goroutines in concurrency to use per child resolvers(default: unlimited)
  # worker_limit: 1000

# Comment or remove this section to skip Apollo Federation support
federation:
  filename: graph/federation.go
  package: graph
  version: 2
  options:
    computed_requires: true

# Where should any generated models go?
model:
  filename: graph/model/models_gen.go
  package: model

  # Optional: Pass in a path to a new gotpl template to use for generating the models
  # model_template: [your/path/model.gotpl]

# Where should the resolver implementations go?
resolver:
  package: graph
  layout: follow-schema # Only other option is "single-file."

  # Only for single-file layout:
  # filename: graph/resolver.go

  # Only for follow-schema layout:
  dir: graph
  filename_template: "{name}.resolvers.go"

  # Optional: turn on to not generate template comments above resolvers
  # omit_template_comment: false
  # Optional: Pass in a path to a new gotpl template to use for generating resolvers
  # resolver_template: [your/path/resolver.gotpl]
  # Optional: turn on to avoid rewriting existing resolver(s) when generating
  # preserve_resolver: false

# Optional: turn on use ` + "`" + `gqlgen:"fieldName"` + "`" + ` tags in your models
# struct_tag: json

# Optional: turn on to use []Thing instead of []*Thing
# omit_slice_element_pointers: false

# Optional: turn on to omit Is<Name>() methods to interface and unions
# omit_interface_checks : true

# Optional: turn on to skip generation of ComplexityRoot struct content and Complexity function
# omit_complexity: false

# Optional: turn on to not generate any file notice comments in generated files
# omit_gqlgen_file_notice: false

# Optional: turn on to exclude the gqlgen version in the generated file notice. No effect if `omit_gqlgen_file_notice` is true.
# omit_gqlgen_version_in_file_notice: false

# Optional: turn on to exclude root models such as Query and Mutation from the generated models file.
# omit_root_models: false

# Optional: turn on to exclude resolver fields from the generated models file.
# omit_resolver_fields: false

# Optional: turn off to make struct-type struct fields not use pointers
# e.g. type Thing struct { FieldA OtherThing } instead of { FieldA *OtherThing }
# struct_fields_always_pointers: true

# Optional: turn off to make resolvers return values instead of pointers for structs
# resolvers_always_return_pointers: true

# Optional: turn on to return pointers instead of values in unmarshalInput
# return_pointers_in_unmarshalinput: false

# Optional: wrap nullable input fields with Omittable
# nullable_input_omittable: true

# Optional: set to speed up generation time by not performing a final validation pass.
# skip_validation: true

# Optional: set to skip running `go mod tidy` when generating server code
# skip_mod_tidy: true

# Optional: if this is set to true, argument directives that
# decorate a field with a null value will still be called.
#
# This enables argumment directives to not just mutate
# argument values but to set them even if they're null.
call_argument_directives_with_null: true

# This enables gql server to use function syntax for execution context
# instead of generating receiver methods of the execution context.
# use_function_syntax_for_execution_context: true

# Optional: set build tags that will be used to load packages
# go_build_tags:
#  - private
#  - enterprise

# Optional: set to modify the initialisms regarded for Go names
# go_initialisms:
#   replace_defaults: false # if true, the default initialisms will get dropped in favor of the new ones instead of being added
#   initialisms: # List of initialisms to for Go names
#     - 'CC'
#     - 'BCC'

# gqlgen will search for any type names in the schema in these go packages
# if they match it will use them, otherwise it will generate them.
autobind:
#  - "{{.}}/graph/model"

# This section declares type mapping between the GraphQL and go type systems
#
# The first line in each type will be used as defaults for resolver arguments and
# modelgen, the others will be allowed when binding to fields. Configure them to
# your liking
models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  # gqlgen provides a default GraphQL UUID convenience wrapper for github.com/google/uuid
  # but you can override this to provide your own GraphQL UUID implementation
  UUID:
    model:
      - github.com/99designs/gqlgen/graphql.UUID

  # The GraphQL spec explicitly states that the Int type is a signed 32-bit
  # integer. Using Go int or int64 to represent it can lead to unexpected
  # behavior, and some GraphQL tools like Apollo Router will fail when
  # communicating numbers that overflow 32-bits.
  #
  # You may choose to use the custom, built-in Int64 scalar to represent 64-bit
  # integers, or ignore the spec and bind Int to graphql.Int / graphql.Int64
  # (the default behavior of gqlgen). This is fine in simple use cases when you
  # do not need to worry about interoperability and only expect small numbers.
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int32
  Int64:
    model:
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
```

Everything has defaults, so add things as you need.

## Inline config with directives

gqlgen ships with some builtin directives that make it a little easier to manage wiring.

To start using them you first need to define them:

```graphql
directive @goModel(
	model: String
	models: [String!]
	forceGenerate: Boolean
) on OBJECT | INPUT_OBJECT | SCALAR | ENUM | INTERFACE | UNION

directive @goField(
	forceResolver: Boolean
	name: String
	omittable: Boolean
	type: String
) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION

directive @goTag(
	key: String!
	value: String
) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION

directive @goExtraField(
	name: String
	type: String!
	overrideTags: String
	description: String
) repeatable on OBJECT | INPUT_OBJECT
```

> Here be dragons
>
> gqlgen doesnt currently support user-configurable directives for SCALAR, ENUM, INTERFACE or UNION. This only works
> for internal directives. You can track the progress [here](https://github.com/99designs/gqlgen/issues/760)

Now you can use these directives when defining types in your schema:

```graphql
type User
	@goModel(model: "github.com/my/app/models.User")
	@goExtraField(name: "Activated", type: "bool") {
	id: ID! @goField(name: "todoId")
	name: String!
		@goField(forceResolver: true)
		@goTag(key: "xorm", value: "-")
		@goTag(key: "yaml")
}

# This make sense when autobind activated.
type Person @goModel(forceGenerate: true) {
	id: ID!
	name: String!
}
```

The builtin directives `goField`, `goModel`, `goTag` and `goExtraField` are automatically registered to `skip_runtime`. Any directives registered as `skip_runtime` will not exposed during introspection and are used during code generation only.

If you have created a new code generation plugin using a directive which does not require runtime execution, the directive will need to be set to `skip_runtime`.

e.g. a custom directive called `constraint` would be set as `skip_runtime` using the following configuration

```yml
# custom directives which are not exposed during introspection. These directives are
# used for code generation only
directives:
  constraint:
    skip_runtime: true
```

## Configuring for a big projects
For a single gqlgen project, what gqlgen config works best really changes at a few key points during the course of that project's normal growth and evolution. For instance:

1. quick proof of concept - `single-file` layout
2. medium-sized team(s) - `follow-schema` layout
3. many teams - separate schema and resolvers into separate packages as [in this example](https://github.com/99designs/gqlgen/tree/master/_examples/large-project-structure/integration/go.mod) (see below)
4. very large type systems (more than 65,000 methods) - `use_function_syntax_for_execution_context`
However, some will instead choose to adopt GraphQL Federation and split into multiple gqlgen instances before one of these growth points is even reached.

Big projects usually divide into a separate domains, grouping all related files and resources under a folder such as:
```
DomainA
- schema
- resolver
- service
DomainB
- schema
- resolver
- service
```

After first generating `resolvers` section you can comment out the entire resolver section of the `config.yaml`, so that resolvers are **not** auto-generated so you can then design any desired resolver architecture.
This idea is from a discussion [https://github.com/99designs/gqlgen/issues/1253](https://github.com/99designs/gqlgen/issues/1253#issuecomment-664448226)
