variable "resource_name" {}
variable "instance_name" {
  default = "my-instance"
}
variable "object" {
  type = object({
    name = optional(string)
  })
}

moved {
  from = aws_instance.a
  to   = aws_instance.b
}

import {
  id = "i-123545"
  to = aws_instance.b
}
