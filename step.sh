#!/usr/bin/env bash
set -euxo pipefail

# download the Bitrise Build Cache CLI
curl -sSfL 'https://raw.githubusercontent.com/bitrise-io/bitrise-build-cache-cli/main/install/installer.sh' | sh -s -- -b /tmp/bin -d

# run the Bitrise Build Cache CLI
/tmp/bin/bitrise-build-cache enable-for gradle

# !!!
# TODO: pass step inputs to the CLI
# !!!
