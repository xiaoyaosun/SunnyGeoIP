FROM golang:1.8
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN make build
CMD ["./bin/geoip"]

