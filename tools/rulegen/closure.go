package main

var closureLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories()`)

var closureProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create {{ .Lang.Name }} library
    closure_js_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = ["@io_bazel_rules_closure//closure/protobuf:jspb"],
        visibility = kwargs.get("visibility"),
        internal_descriptors = [name_pb + "/descriptor.source.bin"],
        suppress = [
            "JSC_LATE_PROVIDE_ERROR",
            "JSC_UNDEFINED_VARIABLE",
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
            "JSC_STRICT_INEXISTENT_PROPERTY",
            "JSC_POSSIBLE_INEXISTENT_PROPERTY",
            "JSC_UNRECOGNIZED_TYPE_ERROR",
        ],
    )`)

func makeClosure() *Language {
	return &Language{
		Dir:   "closure",
		Name:  "closure",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:           "closure_proto_compile",
				Kind:           "proto",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//closure:js"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates closure *.js protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "closure_proto_library",
				Kind:           "proto",
				Implementation: closureProtoLibraryRuleTemplate,
				Usage:          closureLibraryUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates a closure_library with compiled protobuf *.js files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags: []*Flag{
					{
						Category: "build",
						Name:     "incompatible_disallow_struct_provider_syntax",
						Value:    "false",
					},
					{
						Category: "build",
						Name:     "incompatible_use_toolchain_resolution_for_java_rules",
						Value:    "false",
					},
				},
			},
		},
	}
}