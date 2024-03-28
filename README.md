# Activate Bitrise Build Cache Add-On for Gradle

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-activate-gradle-remote-cache?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-activate-gradle-remote-cache/releases)

Activates Bitrise Remote Build Cache add-on for subsequent Gradle builds in the workflow

<details>
<summary>Description</summary> test

This Step activates Bitrise's remote build cache add-on for subsequent Gradle executions in the workflow.

After this Step executes, Gradle builds will automatically read from the remote cache and push new entries if it's enabled.

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `push` | Whether the build can not only read, but write new entries to the remote cache | required | `true` |
| `validation_level` | Level of cache entry validation for both uploads and downloads.  Levels: - `none`: no validation. - `warning`: print a warning about invalid cache entries, but don't interrupt the build - `error`: print an error about invalid cache entries and interrupt the build | required | `warning` |
| `collect_metrics` | When enabled, this sets up Gradle build metrics collection for the subsequent Gradle invocations in the workflow. Metrics are sent to [Bitrise Insights](https://app.bitrise.io/insights). | required | `false` |
| `verbose` | Enable logging additional information for troubleshooting | required | `false` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-activate-gradle-remote-cache/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-activate-gradle-remote-cache/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
