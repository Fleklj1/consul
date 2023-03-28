#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -euo pipefail

upsert_config_entry primary '
kind = "proxy-defaults"
name = "global"
config {
  protocol = "tcp"
}
'

register_services primary

gen_envoy_bootstrap s1 19000 primary
gen_envoy_bootstrap mesh-gateway 19001 primary true
