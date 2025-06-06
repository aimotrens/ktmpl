# Description

ktmpl is a lightweight templating tool to process YAML files. It uses the go template engine.
It was born out of the idea of having a simple way to process Kubernetes manifests in Gitlab CI pipelines.


# Usage

    ktmpl [options] <templates>

&lt;templates&gt; can be a file or a directory.
If a directory is specified, all .yml/.yaml files in that directory will be processed.


# Examples
    ktmpl -i values.yml template.yml          # process template.yml with values.yml
    ktmpl -e -i values.yml template.yml       # process template.yml with values.yml and environment variables
    ktmpl -e -o output.yml template.yml       # process template.yml with environment variables and write to output.yml


# Options
    -v, --version       print version information
    -r, --recursive     recurse into subdirectories
    -o, --output        output file
    -i, --values        values file (as YAML)
    -e, --env           add environment variables to values


# Template functions

(The last argument is the pipeline variable)

    indent(spaces int, s string) string             # indent string with spaces
    substr(start, length int, s string) string      # get substring
    iterate(from, to int) []int                     # create slice of integers
    format(format string, obj any) string           # format string
    toYaml(obj any) string                          # convert object to YAML
    endsWith(suffix, s string) bool                 # check if string ends with suffix
    startsWith(prefix, s string) bool               # check if string starts with prefix
    contains(substring, s string) bool              # check if string contains substring
    include(globPattern string) string              # include files concatinated
    includeAsYamlFields(globPattern string) string  # include files as YAML fields (especially useful for K8s config maps)


# Docker

You can also use ktmpl as a Docker container in you CI/CD pipelines. The image is available on docker hub as `t3a6/kdeploy` and bundles ktmpl and kubectl in one image.


# Example templates and usage

## ./values.yml
```yaml
globalAppName: awesome-app
image: awesome-app:{{.env.CI_COMMIT_REF_SLUG}}
```

## ./template/config.yml
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.globalAppName}}
data:
  APP_ENV: production
```

## ./template/deployment.yml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.globalAppName}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.globalAppName}}
  template:
    metadata:
    labels:
      app: {{.globalAppName}}
    spec:
      containers:
        - name: {{.globalAppName}}
          image: {{.image}}
      envFrom:
        - configMapRef:
            name: {{.globalAppName}}
```

## Usage
```bash
ktmpl -e ./values.yml | ktmpl -i - ./template/ | kubectl apply -f -
```

The first `ktmpl` call will process `values.yml` and replaces the CI environment variable.
The result is piped into the second `ktmpl` call which processes the template files with the values from the first call.
The result from the second call is piped into `kubectl apply` which applies the manifests to the Kubernetes cluster.

The preffered way on Gitlab CI is to write the output to a file and then use `kubectl apply -f xxx.yml` to apply the manifests.
Now you can store the output file as artifact to view the result of the templating process for debugging purposes.


# Exclude directories

To exclude directories from processing, create a `.ktmpl_ignore_dir` file in the directory.
This is especially useful if you have a folder with YAML files that you want to include as config files with includeAsYamlFields.


# Include files as YAML fields

To include files as YAML fields, use the `includeAsYamlFields` function. It accepts a glob pattern as an argument.
If you have previously saved classic config files in K8's config maps, you can now save the config files in clearer separate files and then include them in the config map.

Example:

Your config map:

```yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: my-config-map
    data:
      {{ includeAsYamlFields "config/*" }}
```

Your config files:

```yaml
    # ./config/service-a.yml
    first: value
```

```properties
    # ./config/service-b.properties
    Key1=Value1
```

The `includeAsYamlFields` function will include the contents of `service-a.yml` and `service-b.properties` as YAML fields in the main config map.
The filename is used as the key and the file contents as the value.

Result:

```yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: my-config-map
    data:
      service-a.yml: |
        first: value
      service-b.properties: |
        ConfigA=Value1
```
