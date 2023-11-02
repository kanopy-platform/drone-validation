FROM golang:1.21-alpine

ENV PORT 3000

WORKDIR /go/src/app
COPY . .
RUN go install

CMD ["drone-validation"]
EXPOSE 3000
