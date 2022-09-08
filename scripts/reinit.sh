#!/usr/bin/env sh

# Code generated by shipbuilder init 1.21.0. DO NOT EDIT.

if [ ! -f "./scripts/check.sh" ]; then
  cd $(command dirname -- "$(command readlink -f "$(command -v -- "$0")")")/..
fi

. ./scripts/check.sh

check go jq shipbuilder

# set default type and module
shipbuilder_type="service"
shipbuilder_module="server"

if [ -f "./features/shipbuilder/.env" ]; then
  . ./features/shipbuilder/.env
fi

set -e

export package=$($go mod edit -json | $jq -r .Module.Path)

if [ -n "$shipbuilder_overwrite" ]; then
  for i in $shipbuilder_overwrite; do
    shipbuilder_args="$shipbuilder_args -O $i"
  done
else
  shipbuilder_args="-f"
fi

for i in $shipbuilder_exclude; do
  shipbuilder_args="$shipbuilder_args -x $i"
done
shipbuilder_args="$shipbuilder_args -x version/version_gen.go"

shipbuilder_args="$shipbuilder_args -T $shipbuilder_type"
if [ "$shipbuilder_type" = "service" -a -n "$shipbuilder_module" ]; then
  shipbuilder_args="$shipbuilder_args -M $shipbuilder_module"
fi
shipbuilder_args="$shipbuilder_args $*"
shipbuilder_args="$shipbuilder_args $package"

$shipbuilder init $shipbuilder_args
