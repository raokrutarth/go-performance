
FROM golang:1.12-stretch

EXPOSE 3535
EXPOSE 4000
EXPOSE 5000

# install packages needed for pprof profile pdf generation
RUN apt-get update
RUN apt install -y python-pydot python-pydot-ng graphviz

RUN mkdir /profiles

COPY ./collector/node_exporter/node_exporter /bin
COPY ./collector/process_exporter/process_exporter /bin

COPY ./benchmark-init.sh /


WORKDIR /go/src
CMD ["/bin/bash", "/benchmark-init.sh"]
