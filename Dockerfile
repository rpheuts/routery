FROM golang

RUN go get github.com/rpheuts/routery
RUN go install github.com/rpheuts/routery

# Copy the config file
RUN mkdir /etc/routery
COPY routery.yaml /etc/routery/routery.yaml

WORKDIR /etc/routery
ENTRYPOINT /go/bin/routery

EXPOSE 8080
