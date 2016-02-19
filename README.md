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

WARNING: Please be aware that the exported data format will change in the future!

# Building
Flowbeat uses Glide for dependency management. To install glide see:
https://github.com/Masterminds/glide

"go get github.com/Masterminds/glide" should work in most cases.

Then do

"glide up"

to install all dependencies.
Then build Flowbeat by just running:

"go build"

## Documentation

Set the correct Listen port in flowbeat.yml and start sending it sflow packets.

TODO

## Exported fields

Currently flowbeat exports the raw parsed sflow data. Exactly as it is received and parsed from the wire

