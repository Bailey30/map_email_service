FROM golang:1.21

# create a directory inside the image you are building
WORKDIR /app


# So before you can run go mod download inside your image, you need to get your go.mod and go.sum files copied into it. Use the COPY command to do this.

# In its simplest form, the COPY command takes two parameters. The first parameter tells Docker what files you want to copy into the image. The last parameter tells Docker where you want that file to be copied to.

# Copy the go.mod and go.sum file into your project directory /app which, owing to your use of WORKDIR, is the current directory (./) inside the image. Unlike some modern shells that appear to be indifferent to the use of trailing slash (/), and can figure out what the user meant (most of the time), Docker's COPY command is quite sensitive in its interpretation of the trailing slash.
COPY go.mod go.sum ./

# install go modules
RUN go mod download

# copy souce code into the image
COPY . ./

RUN go build -o /password

EXPOSE 3001

CMD ["/password"]

# to build
# docker build -t password-reset

# to run:
# docker run -d -p 3001:3001 password-reset
# map port 3001 on local machine to port 3001 in container
