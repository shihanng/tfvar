
resource "tfe_variable" "availability_zone_names" {
  key          = "availability_zone_names"
  value        = ["us-west-1a"]
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "docker_ports" {
  key = "docker_ports"
  value = [{
    external = 8300
    internal = 8300
    protocol = "tcp"
  }]
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "image_id" {
  key          = "image_id"
  value        = null
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "password" {
  key          = "password"
  value        = null
  sensitive    = true
  description  = "the root password to use with the database"
  workspace_id = null
  category     = "terraform"
}
