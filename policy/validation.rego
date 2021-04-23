package drone.validation

default deny = false

allowed_types := ["kubernetes"]

is_pipeline {
	input.kind == "pipeline"
}

# drone accepts empty types and use "docker" by default
type = "docker" {
	input.type == ""
} else = input.type {
	true
}

is_valid {
	allowed_types[_] == type
}

deny {
	is_pipeline
	not is_valid
	true
}

out = sprintf("unsupported pipeline types:%v", [type]) {
	deny == true
} else = "" {
	true
}
