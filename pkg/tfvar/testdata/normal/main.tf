variable "resource_name" {}
variable "instance_name" {
  default = "my-instance"
}

moved {
  from = aws_instance.a
  to   = aws_instance.b
}
