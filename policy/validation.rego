package drone.validation

default deny = false

allowed_types := ["kubernetes"]

is_pipeline {
	input.kind == "pipeline"
}

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

out = sprintf("type '%v' is not supported", [type]) {
	deny == true
} else = "" {
	true
}
