FROM node:12.18.3-alpine

RUN apk add --no-cache make pkgconfig gcc g++ python libx11-dev libxkbfile-dev libsecret-dev
RUN apk --no-cache add openjdk11 --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community
ARG version=latest
ENV HOME /home/theia
WORKDIR /home/theia
COPY package.json /home/theia

RUN yarn --pure-lockfile

RUN chmod g+rw /home && \
    mkdir -p /home/project && \
    mkdir -p /home/theia && \
    chown -R node:node /home/theia && \
    chown -R node:node /home/project;

EXPOSE 3000
RUN apk add --no-cache git openssh bash
ENV SHELL=/bin/bash \
    THEIA_DEFAULT_PLUGINS=local-dir:/home/theia/plugins

USER node

ENTRYPOINT [ "node", "/home/theia/src-gen/backend/main.js", "/home/project", "--hostname=0.0.0.0" ]
