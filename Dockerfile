FROM golang:1.18-alpine as build_client
WORKDIR /src
COPY . ./
RUN go build -o pipes-client ./main.go
ENTRYPOINT ["./pipes-client"]

# The --platform flag should only be needed on M1 macbooks (arm64 processors).
FROM --platform=linux/arm64/v8 ubuntu:focal as build_server
WORKDIR /mongod-src
RUN apt update
RUN apt install -y curl
# Download links at https://www.mongodb.com/download-center/community/releases/development
# For arm64 architectures, the Ubuntu 20.04 ARM 64 option works well.
# For x86-64 architectures, the Amazon Linux 2 x64 or RedHat / CentOS 7.0 x64 may work instead.
#
# TODO: The ".deb" and ".rpm" package links on that site are broken. Follow-up
# with in the #server-release Slack channel to see if those can be fixed and we
# can use package managers to do the installation instead.
RUN curl https://fastdl.mongodb.org/linux/mongodb-linux-aarch64-ubuntu2004-6.2.0.tgz -o mongodb620.tgz
RUN tar xvf mongodb620.tgz
RUN rm -r mongodb-data || true
RUN mkdir mongodb-data
RUN rm mongod.log || true
ENTRYPOINT ["mongodb-linux-aarch64-ubuntu2004-6.2.0/bin/mongod"]
CMD ["--dbpath", "mongodb-data", "--logpath", "mongod.log", "--bind_ip", "localhost", "--port", "27017", "--setParameter", "enableComputeMode=true"]
