FROM daocloud.io/golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN make build
CMD ["./bin/geoip"]

