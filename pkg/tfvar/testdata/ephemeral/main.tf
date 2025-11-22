variable "api_token" {
  type        = string
  description = "API token for authentication"
  sensitive   = true
  ephemeral   = true
}

variable "session_id" {
  type      = string
  ephemeral = true
}

variable "regular_var" {
  type    = string
  default = "not-ephemeral"
}

variable "ephemeral_with_default" {
  type      = string
  default   = "temp-value"
  ephemeral = true
}

