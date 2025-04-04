#!/usr/bin/env bash

# 'read' has to be before 'set -e'
read -r -d '' UNAVAILABLE_MESSAGE << EOF_MSG
Bitrise Build Cache is not activated in this build.

You have added the **Activate Bitrise Build Cache for Gradle** add-on step to your workflow.

However, you don't have an activate Bitrise Build Cache Trial or Subscription for the current workspace yet.

You can activate a Trial at [app.bitrise.io/build-cache](https://app.bitrise.io/build-cache),
or contact us at [support@bitrise.io](mailto:support@bitrise.io) to activate it.
EOF_MSG

set -eo pipefail

echo "Checking whether Bitrise Build Cache is activated for this workspace ..."
if [ "$BITRISEIO_BUILD_CACHE_ENABLED" != "true" ]; then
  printf "\n%s\n" "$UNAVAILABLE_MESSAGE"
  set -x
  bitrise plugin install https://github.com/bitrise-io/bitrise-plugins-annotations.git
  bitrise :annotations annotate "$UNAVAILABLE_MESSAGE" --style error || {
    echo "Failed to create annotation"
    exit 3
  }
  exit 2
fi
echo "Bitrise Build Cache is activated in this workspace, configuring the build environment ..."

set -x

# download the Bitrise Build Cache CLI
export BITRISE_BUILD_CACHE_CLI_VERSION="v0.15.26"
curl --retry 5 -sSfL 'https://raw.githubusercontent.com/bitrise-io/bitrise-build-cache-cli/main/install/installer.sh' | sh -s -- -b /tmp/bin -d $BITRISE_BUILD_CACHE_CLI_VERSION

if [ "$collect_metrics" != "true" ] && [ "$collect_metrics" != "false" ]; then
  echo "Parsing inputs failed: Collect Gradle build metrics ($collect_metrics) is not a valid option."
fi

if [ "$push" != "true" ] && [ "$push" != "false" ]; then
  echo "Parsing inputs failed: Push new cache entries ($push) is not a valid option."
fi

if [ "$validation_level" != "none" ] && [ "$validation_level" != "warning" ] && [ "$validation_level" != "error" ]; then
  echo "Parsing inputs failed: Validation level ($validation_level) is not a valid option."
fi

if [ "$verbose" != "true" ] && [ "$verbose" != "false" ]; then
  echo "Parsing inputs failed: Verbose logging ($verbose) is not a valid option."
fi

# run the Bitrise Build Cache CLI
/tmp/bin/bitrise-build-cache enable-for gradle --metrics="$collect_metrics" --push="$push" --validation-level="$validation_level" --debug="$verbose"

