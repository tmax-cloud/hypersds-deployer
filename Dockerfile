FROM ubuntu:18.04

ARG VERSION=15.2.8

RUN apt-get update && apt-get install -y \
    wget software-properties-common lsb-release \
    sshpass 
RUN wget https://download.ceph.com/keys/release.asc &&\
apt-key add release.asc &&\
rm release.asc 
RUN add-apt-repository \
"deb [arch=amd64] https://download.ceph.com/debian-${VERSION}/ \
$(lsb_release -cs) main" &&\
apt-get install -y ceph-common
RUN apt-get purge -y wget software-properties-common lsb-release &&\
apt-get autoremove -y &&\
rm -rf /var/lib/apt/lists/*

RUN mkdir -p /working/config/

WORKDIR /usr/bin/
ADD hypersds-provisioner .
ENTRYPOINT ["hypersds-provisioner"]
