# `tfvar`

[![](https://github.com/shihanng/tfvar/workflows/main/badge.svg?branch=master)](https://github.com/shihanng/tfvar/actions?query=workflow%3Amain)
[![](https://github.com/shihanng/tfvar/workflows/release/badge.svg?branch=master)](https://github.com/shihanng/tfvar/actions?query=workflow%3Arelease)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/shihanng/tfvar)](https://github.com/shihanng/tfvar/releases)
[![Coverage Status](https://coveralls.io/repos/github/shihanng/tfvar/badge.svg?branch=master)](https://coveralls.io/github/shihanng/tfvar?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/shihanng/tfvar)](https://goreportcard.com/report/github.com/shihanng/tfvar)
[![Package Documentation](https://godoc.org/github.com/shihanng/tfvar/pkg/tfvar?status.svg)](http://godoc.org/github.com/shihanng/tfvar/pkg/tfvar)
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

- **tfvar** will search for all input variables and generate template that helps user populates those variables easily:
    ```
    $ tfvar .
    availability_zone_names = ["us-west-1a"]
    docker_ports            = [{ external = 8300, internal = 8300, protocol = "tcp" }]
    image_id                = null
    ```
- Note that default values are assigned to the definitions by default as shown above. Use the `--ignore-default` options to ignore the default values.
    ```
    $ tfvar . --ignore-default
    availability_zone_names = null
    docker_ports            = null
    image_id                = null
    ```
- **tfvar** also provides output in environment variable formats:
    ```
    $ tfvar . -e
    export TF_VAR_availability_zone_names='["us-west-1a"]'
    export TF_VAR_docker_ports='[{ external = 8300, internal = 8300, protocol = "tcp" }]'
    export TF_VAR_image_id=''
    ```
- There is also `--auto-assign` option for those who wants the values from `terraform.tfvars[.json]`, `*.auto.tfvars[.json]`, and environment variables (`TF_VAR_` followed by the name of a declared variable) to be assigned to the generated definitions automatically.
    ```
    $ export TF_VAR_availability_zone_names='["custom_zone"]'
    $ tfvar . --auto-assign
    availability_zone_names = ["custom_zone"]
    docker_ports            = [{ external = 8300, internal = 8300, protocol = "tcp" }]
    image_id                = null
    ```
- Like the [`terraform (plan|apply)`](https://www.terraform.io/docs/configuration/variables.html#variables-on-the-command-line) CLI tool, individual vairables can also be specified via `--var` option.
    ```
    $ tfvar . --var=availability_zone_names='["custom_zone"]' --var=image_id=abc123
    availability_zone_names = ["custom_zone"]
    docker_ports            = [{ external = 8300, internal = 8300, protocol = "tcp" }]
    image_id                = "abc123"
    ```
- Variables in file can also be specified via `--var-file` option.
    ```
    $ cat my.tfvars
    image_id = "xyz"
    $ tfvar . --var-file my.tfvars
    availability_zone_names = ["us-west-1a"]
    docker_ports            = [{ external = 8300, internal = 8300, protocol = "tcp" }]
    image_id                = "xyz"
    ```

- Multiple files can be specified via providing more `--var-file` options, variables overrides as for `terraform` command.
    ```
    $ cat my.tfvars
    image_id = "xyz"

    $ cat other.tfvars
    image_id = "abc"

    $ tfvar . --var-file my.tfvars --var-file other.tfvars
    image_id = "abc"
  ```

For more info, checkout the `--help` page:

```
$ tfvar --help
Generate variable definitions template for Terraform module as
one would write it in variable definitions files (.tfvars).

Usage:
  tfvar [DIR] [flags]

Flags:
  -a, --auto-assign            Use values from environment variables TF_VAR_* and
                               variable definitions files e.g. terraform.tfvars[.json] *.auto.tfvars[.json]
  -d, --debug                  Print debug log on stderr
  -e, --env-var                Print output in export TF_VAR_image_id=ami-abc123 format
  -h, --help                   help for tfvar
      --ignore-default         Do not use defined default values
      --var stringArray        Set a variable in the generated definitions.
                               This flag can be set multiple times.
      --var-file stringArray   Set variables from a file.
                               This flag can be set multiple times.
  -v, --version                version for tfvar
```


## Installation

### [Homebrew (macOS)](https://github.com/shihanng/homebrew-tfvar)

```
brew install shihanng/tfvar/tfvar
```

### Debian, Ubuntu

```
curl -sLO https://github.com/shihanng/tfvar/releases/latest/download/tfvar_linux_amd64.deb
dpkg -i tfvar_linux_amd64.deb
```

### RedHat, CentOS

```
rpm -ivh https://github.com/shihanng/tfvar/releases/latest/download/tfvar_linux_amd64.rpm
```

### Binaries

The [release page](https://github.com/shihanng/tfvar/releases) contains binaries built for various platforms. Download the version matches your environment (e.g. `linux_amd64`) and place the binary in the executable `$PATH` e.g. `/usr/local/bin`:

```
curl -sL https://github.com/shihanng/tfvar/releases/latest/download/tfvar_linux_amd64.tar.gz | \
    tar xz -C /usr/local/bin/ tfvar
```

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

## Contributing

Want to add missing feature? Found bug :bug:? Pull requests and issues are welcome. For major changes, please open an issue first to discuss what you would like to change :heart:.

```
make lint
make test
```

should help with the idiomatic Go styles and unit-tests.

## License
[MIT](./LICENSE)
