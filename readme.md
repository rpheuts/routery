# Routery
___

[![Circle CI](https://img.shields.io/circleci/project/rpheuts/routery.svg)](https://circleci.com/gh/rpheuts/routery)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/rpheuts/routery/blob/master/LICENSE.md)

Routery is aimed to be a convention-over-configuration dynamic reverse proxy for Docker. It can monitor multiple Docker hosts and create dynamic proxy routes for containers based on the container name. The incentive to write this tool came from the frustration of not being able to find a simple tool that would allow dynamic proxying for containers for multiple Docker hosts.

## Features

 - Supports multiple Docker hosts as route providers
 - Allows for multiple frontends and hostnames
 - Supports multiple ports by appending the port name (if not 80) to the host
 - SSL Termination support
 - Basic LDAP Auth support

## How-To

Routery requires only a few bits of information to do its work:
 - Docker host(s)
 - Frontends

Frontends are the ports that Routery will listen on for web requests, you can specify multiple front-ends. In the future it will be possible to assign specific containers to specific domains.

Update the routery.yaml file and update / add frontends and Docker hosts Routery should monitor, it supports both SSL (port 2376) and non-SSL (port 2375) for Docker hosts.
