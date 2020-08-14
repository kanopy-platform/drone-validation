package drone.validation

default deny = false

allowed_types := ["kubernetes"]

type = "docker" {
	input.type == ""
} else = input.type {
	true
}

is_valid {
	allowed_types[_] == type
}

deny {
	not is_valid
	true
}

out = sprintf("type '%v' is not supported", [type]) {
	deny == true
} else = "" {
	true
}
