# Description

ktmpl is a lightweight templating tool to process YAML files.  
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

The last argument is the pipeline variable

    indent(spaces int, s string) string             # indent string with spaces
    substr(start, length int, s string) string      # get substring
    iterate(from, to int) []int                     # create slice of integers
    format(format string, obj any) string           # format string
    toYaml(obj any) string                          # convert object to YAML