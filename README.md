# Flowbeat

Flowbeat is the [Beat](https://www.elastic.co/products/beats) used for
collecting sflow data.

# Current Status
Flowbeat at its current state (2016-02-17) works for router generated Sflow data and imports the following Flow Samples:
 - Basic Link Layer Headers and Protocols (Ethernet, IP, TCP, UDP)
 - ExtendedGateway Flows
 - ExtendedSwitch Flows
 - ExtendedRouter Flows

The used sflow library also parses several Host s-flow samples but this is untested.

Generally if flowbeat does not show the samples you want, the sflow library is probably lacking the parser support for them.

# Building
Flowbeat uses Glide for dependency management. To install glide see:
https://github.com/Masterminds/glide

"go get github.com/Masterminds/glide" should work in most cases.

Then do

"glide up"

to install all dependencies and set 
"export GO15VENDOREXPERIMENT=1" or use Go 1.6

Then build Flowbeat by just running:

"go build"

## Documentation

Set the correct Listen port in flowbeat.yml and start sending it sflow packets.

An elasticsearch mapping is provided which sets all strings to not analyzed and adjust the IP Address fields

TODO

## Exported fields

Currently flowbeat exports the raw parsed sflow data. Exactly as it is received and parsed from the wire

