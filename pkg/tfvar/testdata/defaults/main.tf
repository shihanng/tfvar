variable "instance_name" {
  default = "my-instance"
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

output "availability_zone_names" {
  value = var.availability_zone_names
}

output "instance_name" {
  value = var.instance_name
}
