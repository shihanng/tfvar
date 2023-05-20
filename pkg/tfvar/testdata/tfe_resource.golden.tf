
resource "tfe_variable" "availability_zone_names" {
  key          = "availability_zone_names"
  value        = ["us-west-1a"]
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "aws_amis" {
  key = "aws_amis"
  value = {
    eu-west-1 = "ami-b1cf19c6"
    us-east-1 = "ami-de7ab6b6"
    us-west-1 = "ami-3f75767a"
    us-west-2 = "ami-21f78e11"
  }
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "docker_ports" {
  key = "docker_ports"
  value = [{
    external = 8300
    internal = 8301
    protocol = "tcp"
  }]
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "instance_name" {
  key          = "instance_name"
  value        = "my-instance"
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

resource "tfe_variable" "region" {
  key          = "region"
  value        = null
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "with_optional_attribute" {
  key = "with_optional_attribute"
  value = {
    a = "val-a"
    b = null
    c = 127
  }
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}
