The Github Status plugin allows you to add repo status to a commit from your Drone pipeline.

In it's simplest form, add a status to the build:

```yaml
pipeline:
  github-status:
    when:
      event: pull_request
    image: jmccann/drone-github-status:0.1
    context: ci/custom-status
```

The above will set a status `ci/custom-status` and "inherit" values for `state` (based on `DRONE_BUILD_STATUS`) and `target_url` (based on `DRONE_BUILD_LINK`).

You can set the same `state`, `description` and `target_url` for a list of contexts.  In the following example we set `ci/custom-status-1` and `ci/custom-status-2` to `state` of `pending`.

```diff
pipeline:
  github-status:
    when:
      event: pull_request
    image: jmccann/drone-github-status:0.1
-   context: ci/custom-status
+   context:
+     - ci/custom-status-1
+     - ci/custom-status-2
+   state: pending
```

You can also set status from a value set in a file, defaulting to `failure` if the file does not exist:

```diff
pipeline:
  test:
    image: go:1.10
    commands:
      - go test
      - echo 'success' > .test-status

  github-status:
    when:
      status: [success, failure]
    image: jmccann/drone-github-status:0.1
    context: ci/custom-status
+   file: .test-status
```

# Parameter Reference

`api_key`
: github token to auth to API with

`base_url`
: github API url, defaults to `https://api.github.com`

`context`
: context(s) to create a status for

`description`
: description for the status

`file`
: file to read status from, defaults to `failure` if defined and not found

`state`
: user defined state to set for status.  defaults to drone build status.

`target_url`
: URL for status to link to

`password`
: github password to authenticate with

`username`
: github username to authenticate with