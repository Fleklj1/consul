variable "do_token" {}
variable "key_path" {}
variable "ssh_key_ID" {}
variable "num_instances" {}

## below sourced from 
## https://github.com/hashicorp/terraform/blob/master/examples/digitalocean/variable.tf

# ####
# Current Availiable Datacenter Regions
# As of 05-07-2016
#

variable "do_ams2" {
    description = "Digital Ocean Amsterdam Data Center 2"
    default = "ams2"
}

variable "do_ams3" {
    description = "Digital Ocean Amsterdam Data Center 3"
    default = "ams3"
}

variable "do_fra1" {
    description = "Digital Ocean Frankfurt Data Center 1"
    default = "fra1"
}

variable "do_lon1" {
    description = "Digital Ocean London Data Center 1"
    default = "lon1"
}

variable "do_nyc1" {
    description = "Digital Ocean New York Data Center 1"
    default = "nyc1"
}

variable "do_nyc2" {
    description = "Digital Ocean New York Data Center 2"
    default = "nyc2"
}

variable "do_nyc3" {
    description = "Digital Ocean New York Data Center 3"
    default = "nyc3"
}

variable "do_sfo1" {
    description = "Digital Ocean San Francisco Data Center 1"
    default = "sfo1"
}

variable "do_sgp1" {
    description = "Digital Ocean Singapore Data Center 1"
    default = "sgp1"
}

variable "do_tor1" {
    description = "Digital Ocean Toronto Datacenter 1"
    default = "tor1"
}

# Default Os

variable "ubuntu" {
    description = "Default LTS"
    default = "ubuntu-14-04-x64"
}

variable "centos" {
    description = "Default Centos"
    default = "centos-72-x64"
}

variable "coreos" {
    description = "Defaut Coreos"
    default = "coreos-899.17.0"
}
