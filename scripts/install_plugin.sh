#!/usr/bin/env bash
#
# install_plugins.sh
#
# This script installs the plugin in ~/.terraform.d/plugins

set -e

oss=( linux darwin )
archs=( amd64 386 )
plugins_dir="${HOME}/.terraform.d/plugins"

install_plugin() {
  plugin=$1
  version=0.0.1
  plugin_name=terraform-provider-$(basename "${plugin}")
  plugin_location=$(command -v "${plugin_name}")
  echo "Installing Terraform plugin ${plugin}..."
  for os in "${oss[@]}"
  do
    for arch in "${archs[@]}"
    do
      file="${plugin_name}_v${version}-${os}-${arch}"
      plugin_dst="${plugins_dir}/${plugin}/${version}/${os}_${arch}/${file}"
      mkdir -p "$(dirname "${plugin_dst}")"
      echo "location: ${plugin_location}"
      cp "${plugin_location}" "${plugin_dst}"
      echo "Copied to ${plugin_dst}"
    done
  done
}

install_plugin "terraform.local.com/tapih/slack"