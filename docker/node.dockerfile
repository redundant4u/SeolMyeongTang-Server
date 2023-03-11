FROM node:16-alpine

WORKDIR /home/node

ENV NODE_ENV=prod \
    TZ=Asia/Seoul

RUN apk add --no-cache tzdata build-base python3 openssh-client && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

COPY . .

RUN yarn install --immutable --immutable-cache && yarn build

ENTRYPOINT ["yarn", "node", "dist/main.js"]
