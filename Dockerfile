FROM golang:1.22-alpine AS build-env

WORKDIR /build

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o main-agent code/API/agent/cmd/main.go && go build -o main-orkestrator code/API/orkestrator/cmd/main.go

FROM node:18-alpine

WORKDIR /code/webdaee/

COPY code/webdaee/package.json ./
RUN npm install

COPY code/webdaee/public ./public/
COPY code/webdaee/src ./src/
COPY --from=build-env /build/main-agent ./code/API/agent/cmd/
COPY --from=build-env /build/main-orkestrator ./code/API/orkestrator/cmd/

CMD ["sh", "-c", "ls -la ./code/API/agent/cmd && ls -la ./code/API/orkestrator/cmd && npm run & ./code/API/agent/cmd/main-agent & ./code/API/orkestrator/cmd/main-orkestrator"]
