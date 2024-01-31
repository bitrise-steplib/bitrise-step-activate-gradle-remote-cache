#!/usr/bin/env bash
set -euxo pipefail

# download the Bitrise Build Cache CLI
export BITRISE_BUILD_CACHE_CLI_VERSION="v0.8.0"
curl --retry 5 -sSfL 'https://raw.githubusercontent.com/bitrise-io/bitrise-build-cache-cli/main/install/installer.sh' | sh -s -- -b /tmp/bin -d $BITRISE_BUILD_CACHE_CLI_VERSION

# run the Bitrise Build Cache CLI
/tmp/bin/bitrise-build-cache enable-for gradle

# !!!
# TODO: pass step inputs to the CLI
# !!!
