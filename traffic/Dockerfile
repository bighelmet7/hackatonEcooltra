from golang:alpine

expose 8081
run apk add bash ca-certificates git gcc g++ libc-dev

run mkdir /traffic
workdir /traffic

copy . /traffic
run go build .

cmd ["./traffic"] 
