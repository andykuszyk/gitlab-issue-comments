FROM golang:alpine AS build
WORKDIR /gic
COPY . ./
RUN go build cmd/main.go

FROM alpine
COPY --from=build /gic/main /main
CMD /main
