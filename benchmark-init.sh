#!/bin/bash

# start node exporter to report container level metrics
node_exporter --web.listen-address=":4000" >& /dev/null &

# start process exporter to wait for application "benchmark" to show up
process_exporter -recheck -children=false -web.listen-address ":5000" -procnames "benchmark" >& /dev/null &

dummy() {
    while true; do sleep 6h; done
}

dummy
