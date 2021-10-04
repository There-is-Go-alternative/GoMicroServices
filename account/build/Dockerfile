#######################################################################################################################
# Start from golang base image
FROM golang:1.17 AS builder

# Add Image info
LABEL name='Golang 1.17 (alpine) Image'
LABEL maintainer="Damien Bernard <dbernard.dev@gmail.com>"
LABEL version="1.0"
LABEL description="Temporary Image for Multi stage usage in GoMicroServices"

# Set the current working directory inside the container
WORKDIR /account_service

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download
#RUN go mod vendor

# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/prisma/prisma-client-go prefetch

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# generate the Prisma Client Go client
#RUN go run github.com/prisma/prisma-client-go generate --schema infra/database/schema.prisma
# or, if you use go generate to run the generator, use the following line instead
RUN go generate ./...


#
#FROM setuper as builder
#
## Set the current working directory inside the container
#WORKDIR /account_service

# Build binary
RUN go build -o account .

#######################################################################################################################
# Start from builder stage
FROM builder as runner

# Set the current working directory inside the container
WORKDIR /account_service

# Copy Built binary
COPY --from=builder /account_service/account .

CMD ./account
