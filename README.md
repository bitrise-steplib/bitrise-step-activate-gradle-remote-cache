# Build Cache for Gradle

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-activate-gradle-remote-cache?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-activate-gradle-remote-cache/releases)

Activates Bitrise Remote Build Cache add-on for subsequent Gradle builds in the workflow

<details>
<summary>Description</summary>

This Step activates Bitrise's remote build cache add-on for subsequent Gradle executions in the workflow.

After this Step executes, Gradle builds will automatically read from the remote cache and push new entries if it's enabled.

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `push` | Whether the build can not only read, but write new entries to the remote cache | required | `true` |
| `validation_level` | Level of cache entry validation for both uploads and downloads.  Levels: - `none`: no validation. - `warning`: print a warning about invalid cache entries, but don't interrupt the build - `error`: print an error about invalid cache entries and interrupt the build | required | `warning` |
| `collect_metrics` | When enabled, this sets up Gradle build metrics collection for the subsequent Gradle invocations in the workflow. Metrics are sent to [Bitrise Insights](https://app.bitrise.io/insights). | required | `true` |
| `verbose` | Enable logging additional information for troubleshooting | required | `false` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-activate-gradle-remote-cache/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-activate-gradle-remote-cache/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
