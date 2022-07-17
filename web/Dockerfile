FROM node:16-alpine

ENV VUE_APP_PROXY_URL=http://localhost:8080
ENV VUE_APP_SERVICE_URL=http://localhost:8080/api
ENV VUE_APP_PHONE_NUMBER=6282299921720

RUN mkdir /app

COPY . /app

RUN yarn --cwd /app install

CMD ["yarn", "--cwd", "/app", "serve"]