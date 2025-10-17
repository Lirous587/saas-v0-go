FROM golang:1.25.1
WORKDIR /app
COPY ../saas .
COPY ../.env.docker .env
COPY ../model.conf model.conf
RUN chmod +x saas
CMD ["./saas"]

#$env:GOOS="linux"
#$env:GOARCH="amd64"
#go build -o saas .
