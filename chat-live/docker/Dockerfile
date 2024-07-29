FROM node:18-alpine

WORKDIR /home/node

ENV NODE_ENV=production \
    TZ=Asia/Seoul

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    npm i -g pnpm pm2 typescript

COPY . .

RUN pnpm i && pnpm run build

ENTRYPOINT ["pm2-runtime", "dist/index.js"]
