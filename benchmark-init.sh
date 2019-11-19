#!/bin/bash

# start node exporter to report container level metrics
/exporters/node_exporter --web.listen-address=":4000" &

# start process exporter to wait for application "benchmark" to show up
/exporters/process_exporter -children=false -web.listen-address ":5000" -procnames benchmark &

dummy() {
    while true; do sleep 6h; done
}

dummy
