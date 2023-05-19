variable "resource_name" {}
variable "instance_name" {
  default = "my-instance"
}
variable "object" {
  type = object({
    name = string
  })
}

moved {
  from = aws_instance.a
  to   = aws_instance.b
}
