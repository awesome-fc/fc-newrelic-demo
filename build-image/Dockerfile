FROM golang:1.12.16-stretch
RUN mkdir -p $GOPATH/src/golang.org/x/
RUN cd $GOPATH/src/golang.org/x/ && git clone https://github.com/golang/net.git
RUN cd $GOPATH/src/golang.org/x/ &&  git clone https://github.com/golang/sys.git
RUN go get github.com/newrelic/go-agent

#RUN mkdir -p $GOPATH/src/google.golang.org/
#RUN cd $GOPATH/src/google.golang.org/ && git clone https://github.com/grpc/grpc-go
##./google.golang.org/grpc
#RUN cd $GOPATH/src/google.golang.org/ && git clone https://github.com/google/go-genproto.git
##./google.golang.org/genproto
#RUN cd $GOPATH/src/google.golang.org/ && git clone https://github.com/protocolbuffers/protobuf-go
##./google.golang.org/protobuf/


