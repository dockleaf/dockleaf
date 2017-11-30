#!/bin/sh

docker run -e DOCKLEAF_DEFINITION='/config/definition.json' -e DOCKLEAF_VERSION='/config/version.json' -v `pwd`:/config dockleaf:latest