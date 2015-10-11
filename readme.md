# Routery
___

[![Circle CI](https://img.shields.io/circleci/project/rpheuts/routery.svg)](https://circleci.com/gh/rpheuts/routery)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/rpheuts/routery/blob/master/LICENSE.md)

Routery is aimed to be a convention-over-configuration dynamic reverse proxy for Docker. It can monitor multiple Docker hosts and create dynamic proxy routes for containers based on the container name.

## Features

 - Supports multiple Docker hosts as router providers
 - Allows for multiple frontends and hostnames
 - Supports multiple ports by appending the port name (if not 80) to the host
