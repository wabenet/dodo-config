package spec

backdrops: [string]: #Backdrop

#Backdrop: {
  name?:           string
  aliases?:        [...string]
  container_name?: =~"[a-zA-Z0-9][a-zA-Z0-9_.-]*"
  image?:          string | #BuildInfo
  build?:          #BuildInfo
  runtime?:        string
  script?:         string
  interpreter?:    [...string]
  user?:           string
  working_dir?:    string
  environment:     #Environment | [...#EnvironmentVariable] | [...string] | *[]
  ports:           #Ports       | [...#PortMapping]         | [...string] | *[]
  mounts:          #Mounts      | [...#Mount]               | [...string] | *[]
  capabilities:    [...string]  | *[]

  // Deprecated
  volumes:         #Volumes     | [...#VolumeMount]         | [...string] | *[]
  devices:         #Devices     | [...#DeviceMapping]       | [...string] | *[]

  ...
}

#BuildInfo: {
  name?:         string
  builder?:      string
  dependencies:  [...string] | *[]
  context?:      string
  dockerfile?:   string
  steps?:        string
  arguments:     #BuildArguments | [...#BuildArgument] | *[]
  secrets:       #BuildSecrets   | [...#BuildSecret]   | *[]
  ssh_agents:    #SSHAgents      | [...#SSHAgent]      | *[]
  ...
}

#Environment: [string]: #EnvironmentVariable

#EnvironmentVariable: {
  name?:  string
  value?: string
}

#Ports: [string]: #PortMapping

#PortMapping: {
  target:    string | int
  publish:   string | int
  protocol?: string
  host_ip?:  string
}

// Deprecated
#Volumes: [string]: #VolumeMount | *#BindMount

// Deprecated
#Devices: [string]: #DeviceMapping

// Deprecated
#DeviceMapping: #DeviceMount | #DeviceRule

#Mounts: [string]: #Mount

#Mount: #BindMount | #VolumeMount | #TmpfsMount | #ImageMount | #DeviceMount | #DeviceRule

#BindMount: {
  type:     "bind"
  source:   =~"/.*"
  target?:  string
  readonly: bool | *false
}

#VolumeMount: {
  type:     "volume"
  source:   =~"[^/].*"
  target:   string
  path?:    string
  readonly: bool | *false
}

#TmpfsMount: {
  type:  "tmpfs"
  path:  string
  size?: int
  mode?: string
}

#ImageMount: {
  type:     "image"
  source:   string
  target:   string
  path?:    string
  readonly: bool | *false
}

#DeviceMount: {
  type:        "device"
  source:       string
  target?:      string
  permissions?: string
}

#DeviceRule: {
  type:        "device"
  cgroup_rule: string
}

#BuildArguments: [string]: #BuildArgument

#BuildArgument: {
  name:  string
  value: string
}

#BuildSecrets: [string]: #BuildSecret

#BuildSecret: {
  id:   string
  path: string
}

#SSHAgents: [string]: #SSHAgent

#SSHAgent: {
  id:            string
  identity_file: string
}
