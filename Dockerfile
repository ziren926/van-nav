# Stage 1: Build frontend
FROM node:18-alpine AS frontendbuilder
WORKDIR /app
COPY . .
RUN npm install -g pnpm
RUN cd /app && cd ui && pnpm install && CI=false pnpm build && cd ..
RUN cd /app && mkdir -p public
RUN cp -r ui/build/* public/

# Stage 2: Build backend
FROM golang:1.19-alpine3.18 AS binarybuilder
RUN apk --no-cache --no-progress add  git
WORKDIR /app
COPY . .
COPY --from=frontendbuilder /app/public /app/public
# 删除 pnpm install .，只保留 Go 相关的命令
RUN cd /app && ls -la && go mod tidy && go build -o nav

# Stage 3: Final image
FROM alpine:latest
ENV TZ="Asia/Shanghai"
RUN apk --no-cache --no-progress add \
    ca-certificates \
    tzdata && \
    cp "/usr/share/zoneinfo/$TZ" /etc/localtime && \
    echo "$TZ" >  /etc/timezone
WORKDIR /app
COPY --from=binarybuilder /app/nav /app/

VOLUME ["/app/data"]
EXPOSE 6412
ENTRYPOINT [ "/app/nav" ]