# `tfvar`

[![](https://github.com/shihanng/tfvar/workflows/main/badge.svg?branch=master)](https://github.com/shihanng/tfvar/actions?query=workflow%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/shihanng/tfvar/badge.svg?branch=master)](https://coveralls.io/github/shihanng/tfvar?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/shihanng/tfvar)](https://goreportcard.com/report/github.com/shihanng/tfvar)
[![GitHub license](https://img.shields.io/github/license/shihanng/tfvar)](https://github.com/shihanng/tfvar/blob/master/LICENSE)

**tfvar** is a [Terraform](https://www.terraform.io/)'s [variable definitions](https://www.terraform.io/docs/configuration/variables.html#assigning-values-to-root-module-variables) template generator.

For Terraform configuration that has input variables declared, e.g.,

```terraform
variable "image_id" {
  type = string
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}
```

**tfvar** will search for all input variables and generate template that helps user populates those variables easily:

```
$ tfvar .
availability_zone_names = ["us-west-1a"]
docker_ports            = [{ external = 8300, internal = 8300, protocol = "tcp" }]
image_id                = null
```

## Installation

### For Gophers

With [Go](https://golang.org/doc/install) already installed in your system, use `go get`

```
go get github.com/shihanng/tfvar
```

or clone this repo and `make install`

```
git clone https://github.com/shihanng/tfvar.git
cd tfvar
make install
```

## Usage

```
$ tfvar --help
Generate variable definitions template for Terraform module as
one would write it in variable definitions files (.tfvars).

Usage:
  tfvar [DIR] [flags]

Flags:
  -a, --auto-assign      Use values from environment variables TF_VAR_* and
                         variable definitions files e.g. terraform.tfvars[.json] *.auto.tfvars[.json]
  -d, --debug            Print debug log on stderr
  -e, --env-var          Print output in export TF_VAR_image_id=ami-abc123 format
  -h, --help             help for tfvar
      --ignore-default   Do not use defined default values
```

## Contributing

Want to add missing feature? Found bug :bug:? Pull requests and issues are welcome. For major changes, please open an issue first to discuss what you would like to change :heart:.

```
make lint
make test
```

should help with the idiomatic Go styles and unit-tests.

## License
[MIT](./LICENSE)
