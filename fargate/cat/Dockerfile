FROM golang:1.14-buster as build
RUN mkdir /app
ADD . /app
WORKDIR /app
ENV GOPROXY=direct
RUN go build -o cat .

FROM gcr.io/distroless/base-debian10
COPY --from=build /app/cat /

CMD ["./cat"]

EXPOSE 8080
