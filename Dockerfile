# syntax=docker/dockerfile:1

FROM golang:1.19-alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY / ./
RUN ls -la ./*.go

RUN go build -o /thermo_status_server

EXPOSE 8080

CMD [ "/thermo_status_server" ]


#****************** To Build and Deploy to DockerHub**************************

#docker login
#docker buildx build --platform=linux/arm64,linux/amd64 -t speshl/thermo_status_server:latest --push .

#***************************************************************************


#****************** Build and Test Local**************************
#docker build --platform=linux/amd64 -t speshl/thermo_status_server:latest .

#docker run --env-file pi.env -d -p 8080:8080 speshl/thermo_status_server:latest



#*******************Push to docker Hub OR Save/Load locally ******************************
#docker push speshl/thermo_status_server:tagname

#docker save --output thermo_status_server.tar thermo_status_server
#docker load --input thermo_status_server.tar