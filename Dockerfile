FROM golang:alpine AS goaccess-api

# Working directory for build
WORKDIR /build

# Download dependecies from go.mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build the code to get the binary (goaccess)
COPY . .
WORKDIR /build/cmd/goaccess-api/
RUN go get
RUN go build -o goaccessapi

# # Working directory for place the binary
WORKDIR /bin
RUN mkdir init

RUN cp /build/cmd/goaccess-api/goaccessapi .
RUN cp -R /build/init/modules/ ./init
EXPOSE 8077
ENTRYPOINT ["/bin/goaccessapi"]
# RUN ping google.com