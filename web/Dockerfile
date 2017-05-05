FROM node:7.10.0-alpine

COPY package.json /web/package.json
COPY yarn.lock /web/yarn.lock

WORKDIR /web

RUN yarn install

COPY . /web/

CMD ["yarn", "run", "build:prod"]
