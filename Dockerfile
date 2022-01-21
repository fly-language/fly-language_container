FROM node:12.18.3-alpine

LABEL version="0.0.1"

# install dependencies of the container
RUN apk add --no-cache make pkgconfig gcc g++ python libx11-dev libxkbfile-dev libsecret-dev
RUN apk --no-cache add openjdk11 --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community
ARG version=latest
ENV HOME /home/theia
WORKDIR /home/theia
COPY package.json /home/theia

# generate the theia IDE
RUN yarn --pure-lockfile

# create folders 
RUN chmod g+rw /home && \
    mkdir -p /home/project && \
    mkdir -p /home/theia && \
    mkdir -p /home/fly/libcle && \
    chown -R node:node /home/theia && \
    chown -R node:node /home/project && \
    chown -R node:node /home/fly/lib;

# copy fly dependencies 
COPY ./lib /home/fly/lib

# expose port and set default shell stuff
EXPOSE 3000
RUN apk add --no-cache git openssh bash
ENV SHELL=/bin/bash \
    THEIA_DEFAULT_PLUGINS=local-dir:/home/theia/plugins

USER node

ENTRYPOINT [ "node", "/home/theia/src-gen/backend/main.js", "/home/project", "--hostname=0.0.0.0" ]
