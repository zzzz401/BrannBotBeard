FROM golang:1.7

MAINTAINER Aidan Law afl@aidan-law.com

ADD Brann.go /go/src/Brann/Brann.go

ADD . /go/bin/

RUN go get github.com/bwmarrin/discordgo

RUN go install Brann

WORKDIR /go/bin

#Add Discord Token

CMD ./Brann -t
