FROM golang

ENV PACKAGE_PATH /go/src/github.com/metal-go/metal
RUN mkdir -p $PACKAGE_PATH
WORKDIR $PACKAGE_PATH
ADD docker/builds/rexecd/cmd /
RUN chmod +x /cmd
CMD ["/cmd"]
