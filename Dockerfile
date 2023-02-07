FROM --platform=linux/arm64/v8 ubuntu:focal as build_server
WORKDIR /src

# 1. Copy over the Go files and Go binary.
COPY . ./

# 2. Install mongod.
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
RUN tar -xzf mongodb620.tgz
RUN rm -r mongodb-data || true
RUN mkdir mongodb-data
RUN rm mongod.log || true

ENTRYPOINT ["./pipes-client-linux-arm64"]
