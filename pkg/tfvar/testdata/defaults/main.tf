variable "region" {}

variable "instance_name" {
  default = "my-instance"
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

output "region" {
  value = var.region
}

output "instance_name" {
  value = var.instance_name
}

output "availability_zone_names" {
  value = var.availability_zone_names
}

variable "aws_amis" {
  default = {
    "eu-west-1" = "ami-b1cf19c6"
    "us-east-1" = "ami-de7ab6b6"
    "us-west-1" = "ami-3f75767a"
    "us-west-2" = "ami-21f78e11"
  }
}

variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8301
      external = 8300
      protocol = "tcp"
    }
  ]
}

variable "password" {
  type        = string
  description = "the root password to use with the database"
  sensitive   = true
}
