#######################################################################################################################
# Start from golang base image
FROM golang:1.17-alpine AS builder

# Add Image info
LABEL name='Golang 1.17 (alpine) Image'
LABEL maintainer="Okil Saber LAKHDARI <lakhsaber@gmail.com>"
LABEL version="1.0"
LABEL description="Temporary Image for Multi stage usage in GoMicroServices"

# Set the current working directory inside the container
WORKDIR /chats_service

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download
#RUN go mod vendor

# Copy the source from the current directory to the working Directory inside the container
COPY . .

RUN go build -o chats .

#######################################################################################################################
# Start from builder stage
FROM builder as runner

# Set the current working directory inside the container
WORKDIR /chats_service

# Copy Built binary
COPY --from=builder /chats_service/chats .

CMD ./chats
