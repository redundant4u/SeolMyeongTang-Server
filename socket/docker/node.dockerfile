FROM node:16-alpine

WORKDIR /home/node

ENV NODE_ENV=production \
    TZ=Asia/Seoul

RUN apk add --no-cache tzdata build-base python3 openssh-client && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    npm i -g pnpm pm2

COPY . .

RUN pnpm i

ENTRYPOINT ["pm2-runtime", "dist/main.js"]
