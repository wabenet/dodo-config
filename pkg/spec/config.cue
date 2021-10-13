package spec

include: [...#Include] | *[]

#Include: {
  file: string
}

backdrops: [string]: #Backdrop

#Backdrop: {
  name?:           string
  aliases?:        [...string]
  container_name?: =~"[a-zA-Z0-9][a-zA-Z0-9_.-]*"
  image?:          string | #BuildInfo
  build?:          #BuildInfo
  runtime?:        string
  script?:         string
  interpreter:     [...string] | *["/bin/sh", "-c"]
  user?:           string
  working_dir?:    string
  environment:     #Environment | [...#EnvironmentVariable] | [...string] | *[]
  ports:           #Ports       | [...#PortMapping]         | [...string] | *[]
  volumes:         #Volumes     | [...#VolumeMount]         | [...string] | *[]
  devices:         #Devices     | [...#DeviceMapping]       | [...string] | *[]
  capabilities:    [...string] | *[]
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

#Volumes: [string]: #VolumeMount

#VolumeMount: {
  source:   string
  target?:  string
  readonly: bool | *false
}

#Devices: [string]: #DeviceMapping

#DeviceMapping: #DeviceMount | #DeviceRule

#DeviceMount: {
  source:       string
  target?:      string
  permissions?: string
}

#DeviceRule: {
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
