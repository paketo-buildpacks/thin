# Thin Paketo Cloud Native Buildpack

## `gcr.io/paketo-buildpacks/thin`

The Thin CNB sets the start command for a given ruby application that runs on a Thin server.

## Configuration

### Thin config file

It is possible to provide a thin config file to be passed to thin as part of
the start command (via the `thin -C <some-file>` option).

The `BP_THIN_CONFIG_LOCATION` environment variable allows you to specify the
location of a thin config file. This can either be an absolute path, or a
relative path (relative to the application root directory).

If `BP_THIN_CONFIG_LOCATION` is unset and there is a `thin.yml` file in the
application root directory, this file will be provided to thin as part of the
start command.

If `BP_THIN_CONFIG_LOCATION` is set to a value that does not correspond to a
file, the build phase will fail.

### `buildpack.yml` Configurations

There are no extra configurations for this buildpack based on `buildpack.yml`.

## Integration

This CNB writes a start command, so there's currently no scenario we can
imagine that you would need to require it as dependency. If a user likes to
include some other functionality, it can be done independent of the Thin CNB
without requiring a dependency of it.

To package this buildpack for consumption:
```
$ ./scripts/package.sh
```
This builds the buildpack's source using GOOS=linux by default. You can supply another value as the first argument to package.sh.
