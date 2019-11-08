#!/bin/bash

# start process exporter to wait for application "benchmark" to show up
/exporters/process_exporter -web.listen-address ":5000" -procnames benchmark &

# start node exporter to report container level metrics
/exporters/node_exporter --web.listen-address=":4000" &

dummy() {
    while true; do sleep 30m; done
}

dummy
