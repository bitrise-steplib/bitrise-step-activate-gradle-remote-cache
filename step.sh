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
export BITRISE_BUILD_CACHE_CLI_VERSION="v0.17.10"
curl --retry 5 -m 30 -sSfL 'https://raw.githubusercontent.com/bitrise-io/bitrise-build-cache-cli/main/install/installer.sh' | sh -s -- -b /tmp/bin -d $BITRISE_BUILD_CACHE_CLI_VERSION || true

# Fall back to Artifact Registry if the download failed
if [ ! -f /tmp/bin/bitrise-build-cache ]; then
  echo "Failed to download Bitrise Build Cache CLI, trying Artifact Registry ..."

  version="${BITRISE_BUILD_CACHE_CLI_VERSION#v}"
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  arch=$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')
  package="bitrise-build-cache_${os}_${arch}.tar.gz"
  filename="bitrise-build-cache_${version}_${os}_${arch}.tar.gz"

  filepath="$package:$version:$filename"

  echo "Downloading Bitrise Build Cache CLI from Artifact Registry: ${filepath}"

  curl --retry 5 -m 60 -sSfL "https://artifactregistry.googleapis.com/download/v1/projects/ip-build-cache-prod/locations/us-central1/repositories/build-cache-cli-releases/files/${filepath}:download?alt=media" -o $package
  tar -xzf "$package"
  mkdir -p /tmp/bin
  mv "bitrise-build-cache" /tmp/bin/bitrise-build-cache
  rm -rf "$package"
fi

if [ ! -f /tmp/bin/bitrise-build-cache ]; then
  echo "Failed to download Bitrise Build Cache CLI, exiting."
  exit 1
fi

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
/tmp/bin/bitrise-build-cache activate gradle --debug="$verbose" --cache --cache-push="$push" --cache-validation="$validation_level" --analytics="$collect_metrics" 
