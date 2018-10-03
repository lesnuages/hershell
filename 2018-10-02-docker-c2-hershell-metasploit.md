Using Docker to host a hershell and Metasploit C2
===

# Goal

We want to use Docker to build and fresh hershell implant and easily distribute it.  

Hershell comes with the ability to upgrade the infected target to a `meterpreter` implant. We'll also be using Docker to host our second stage C2 (C3 in this article)

Once done, generating, sharing and managing the implant should be a fast and cloud-native process.

# Requirements

* `docker` environment
* Targets (use windows isos)
* Some place to host your C2 and C3
    * AWS or GCE free tier
    * local docker instance

# Why hershell

* golang
* minimal
* session upgrade

# Steps

## Creating a Hershell build with a Dockerfile

### Building

* golang + alpine
* golang image bin folder
* packing
* moving the certs

### Sharing

* go serve, like python simplehttpserver but in golang

## Summoning Metasploit in the cloud

### Launch

* remnux
* postgresql
* open port 8443 (-p "8443:8443")

### Use


* launch msfconsole
* set params (make a rc file?)
* upgrade


## More tricks

### Using docker engine

* switch easily between instances and providers

### Make a rc file

* following input to msf from hershell readme