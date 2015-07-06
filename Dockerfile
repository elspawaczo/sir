FROM gliderlabs/alpine:3.1

ENV GOPATH /sir
ENV GOBIN $GOPATH/bin
ENV PATH $PATH:$GOBIN

RUN apk update && apk add curl git build-base bash mercurial bzr go && rm -rf /var/cache/apk/*

RUN git clone https://github.com/pote/gpm.git \
    && cd gpm \
    && git checkout v1.3.2 \
    && ./configure \
    && make install

WORKDIR /sir

COPY ./Godeps /sir/Godeps
RUN gpm install

COPY . /sir/src/github.com/thisissoon/sir

WORKDIR /sir/src/github.com/thisissoon/sir

RUN go install cmd/sir/sir.go

EXPOSE 8000

CMD ["sir"]
