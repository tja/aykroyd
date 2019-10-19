ARG BASEIMAGE=alpine:3.10


#
# Build "aykroyd" binary
#
FROM $BASEIMAGE AS build

RUN apk --no-cache add build-base \
                       git \
                       go

WORKDIR /go/src/app
COPY . .

ENV GO111MODULE=on
RUN go build -ldflags="-s -w" -o aykroyd main.go


#
# Bundle "aykroyd" into production image
#
FROM $BASEIMAGE AS prod

LABEL maintainer="thomas@crissyfield.de"
LABEL version="1.0.0"
LABEL description="Email forwards via Postfix"

RUN apk --no-cache add tzdata

WORKDIR /app
COPY --from=build /go/src/app/aykroyd /app

EXPOSE 80

CMD [ "/app/aykroyd" ]
