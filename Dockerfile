FROM golang:1.18-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o ec2 .

CMD ["/app/ec2"]