#!/usr/bin/env bash
set -eo pipefail

UNAVAILABLE_MESSAGE=$(cat <<-END
Bitrise Build Cache is not activated in this build.

You have added the **Activate Bitrise Build Cache for Gradle** add-on step to your workflow.

However, you don't have an activate Bitrise Build Cache Trial or Subscription for the current workspace yet.

You can activate a Trial at [app.bitrise.io/build-cache](https://app.bitrise.io/build-cache),
or contact us at [support@bitrise.io](mailto:support@bitrise.io) to activate it.
END
)

if [ "$BITRISEIO_BUILD_CACHE_ENABLED" != "true" ]; then
  printf "\n%s\n" "$UNAVAILABLE_MESSAGE"
  set -x
  bitrise plugin install https://github.com/bitrise-io/bitrise-plugins-annotations.git
  bitrise :annotations annotate "$UNAVAILABLE_MESSAGE" --style error || {
    echo "Failed to create annotation"
    exit 1
  }
  exit 1
fi

set -x

# download the Bitrise Build Cache CLI
export BITRISE_BUILD_CACHE_CLI_VERSION="v0.8.0"
curl --retry 5 -sSfL 'https://raw.githubusercontent.com/bitrise-io/bitrise-build-cache-cli/main/install/installer.sh' | sh -s -- -b /tmp/bin -d $BITRISE_BUILD_CACHE_CLI_VERSION

# run the Bitrise Build Cache CLI
/tmp/bin/bitrise-build-cache enable-for gradle

# !!!
# TODO: pass step inputs to the CLI
# !!!
