FROM golang:1.22-alpine AS build-env

WORKDIR /build

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o main-agent cmd/agent/main.go && go build -o main-orkestrator cmd/orkestrator/main.go

FROM node:18-alpine

WORKDIR build/

COPY webdaee/package.json ./
RUN npm install

COPY webdaee/public ./public/
COPY webdaee/src ./src/
COPY --from=build-env /build/main-agent ./code/API/agent/cmd/
COPY --from=build-env /build/main-orkestrator ./code/API/orkestrator/cmd/

CMD ["sh", "-c", "ls -la ./cmd/agent && ls -la ./cmd/orkestrator && npm start & ./code/API/agent/cmd/main-agent & ./code/API/orkestrator/cmd/main-orkestrator"]
