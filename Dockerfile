FROM node:16-alpine as build-stage
WORKDIR /app
COPY lnkd-front/package*.json ./
RUN npm install
COPY ./lnkd-front/ .
RUN npm run build

FROM golang:1.18-alpine
WORKDIR /app
COPY --from=build-stage /app/build lnkd-front/build/
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN go get -t .

RUN go build -o /lnkd

EXPOSE 8070
CMD ["/lnkd"]
