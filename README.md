# dodo config

Adds support to load dodo backdrops from configuration files.

## installation

This plugin is already included in the dodo default distribution.

If you want to compile your own dodo distribution, you can add this plugin with the
following generate config:

```yaml
plugins:
  - import: github.com/wabenet/dodo-config/pkg/plugin
```

Alternatively, you can install it as a standalone plugin by downloading the
correct file for your system from the [releases page](https://github.com/wabenet/dodo-config/releases),
then copy it into the dodo plugin directory (`${HOME}/.dodo/plugins`).

## configuration

The plugin searches for configuration files written in YAML or JSON. Search paths
include the current directory, any parent directory up to the filesystem root,
as well as all the usual places for config files (`$HOME`, `$XDG_CONFIG_HOME/dodo`,
`%APPDATA%`, etc). The configuration file name must be some variation of
`dodo.yml`. Either with a `.yml` or `.yaml` or even `.json` extension, and an
optional leading dot.

The exact syntax of the configuration file is specified via the [CUE](https://cuelang.org/)
constraints [located here](pkg/spec/config.cue).

Additionally, all fields are templated via the go templating engine.
See the [templating](#templating)section for details.

### backdrops

Backdrops are the main components of dodo. Each backdrop is a template for a
container that acts as runtime environment for a script. The top-level configuration
object in each file is `backdrops`, which is a map of backdrop names to objects
with the following options:

* `aliases`: a list of aliases that can be used instead of the backdrop name to run it
* `container_name`: set the container name
* `image` or `build`: defines configuration for building the container image. Can
  be either a string containing an existing image name, or an object with:
  * `name`: a name (tag) for the resulting image
  * `builder`: name of the build plugin to use, select the first one available by default
  * `context`: path to the build context
  * `dockerfile`: path to the dockerfile (relative to the build context)
  * `steps`: Alternatively, supplies an inline dockerfile
  * `arguments`: build arguments for the image
  * `secrets`: secrets used for building
  * `ssh_agents`: ssh agent connfiguration used for building
  * `dependencies`: list of image names that are required to build this image.
    Other backdrop configurations are searched for build declarations with this
    `name` field, and are built before this image.
* `environment`: set environment variables
* `volumes`: list of additional bind-mount volumes
* `ports`: list of exposed ports from the container
* `user`: set the uid inside the container
* `working_dir`: set the working directory inside the container
* `script`: the script that should be executed
* `interpreter`: set the interpreter that should execute the script (defaults to
  `/bin/sh`)

### templating

All strings in the YAML configuration are processed by the [golang templating
engine](https://golang.org/pkg/text/template/). The following additional methods
are available:

 * Everything from the [sprig library](http://masterminds.github.io/sprig/)
 * `{{ cwd }}` evaluates to the current working directory
 * `{{ currentFile }}` is the path to the current YAML file that is evaluated
 * `{{ currentDir }}` is the path to the directory where the current file is located
 * `{{ projectRoot }}` is the path to the current Git project (determined by the
   first `.git` directory found by walking the current working directory upwards).
   Useful in combination with `{{ projectPath }]` if you don't only want to
   bind-mount the current directory but the whole project.
 * `{{ projectPath }}` the path of the current working directory relative to `{{
   projectRoot }}`
 * `{{ env <variable> }}` evaluates to the contents of environment variable
   `<variable>`
 * `{{ user }}` evaluates to the current user, in form of a
   [golang user](https://golang.org/pkg/os/user/). From this, you can access
   fields like `{{ user.HomeDir }}` or `{{ user.Uid }}`.
 * `{{ sh <command> }}` executes `<command>` via `/bin/sh` and evaluates
   to its stdout

### includes

Includes allow merging additional files that are not in the search path into the
current configuration. For example:

```
include:
  - file: '{{ currentDir }}/some_other_file.yaml'
```


## license & authors

```text
Copyright 2022 Ole Claussen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
