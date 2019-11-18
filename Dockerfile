
FROM golang:1.12-stretch

EXPOSE 3535
EXPOSE 4000
EXPOSE 5000

# install packages needed for pprof profile pdf generation
RUN apt-get update && apt install python-pydot python-pydot-ng graphviz

RUN mkdir /exporters
COPY ./collector/node_exporter/node_exporter /exporters
COPY ./collector/process_exporter/process_exporter /exporters/

COPY ./benchmark-init.sh /


WORKDIR /go/src
CMD ["/bin/bash", "/benchmark-init.sh"]
