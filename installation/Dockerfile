FROM node:12.18.3-alpine

LABEL version="debug"

# install dependencies of the container
RUN apk add --no-cache make pkgconfig gcc g++ python libx11-dev libxkbfile-dev libsecret-dev zip gradle maven python-dev
RUN apk --no-cache add openjdk11 --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community
ARG version=latest
ENV HOME /home/theia
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/home/theia/.local/bin/
WORKDIR /home/theia

RUN pip3 install --upgrade pip && \
    pip3 install --no-cache-dir awscli && \
    pip3 install --no-cache-dir virtualenv

COPY package.json /home/theia

# generate the theia IDE
RUN yarn --pure-lockfile

# create folders
RUN chmod g+rw /home && \
    mkdir -p /home/project && \
    mkdir -p /home/theia && \
    mkdir -p /home/fly && \
    chown -R node:node /home/theia && \
    chown -R node:node /home/project && \
    chown -R node:node /home/fly;

# expose port and set default shell stuff
EXPOSE 3000
RUN apk add --no-cache git openssh bash
ENV SHELL=/bin/bash \
    THEIA_DEFAULT_PLUGINS=local-dir:/home/theia/plugins

# set default user
USER node

# install dependecies to run fly
RUN mkdir -p /home/fly/src/main/fly
WORKDIR /home/fly
COPY pom.xml /home/fly/
RUN mvn install

WORKDIR /home/theia

ENTRYPOINT [ "node", "/home/theia/src-gen/backend/main.js", "/home/project", "--hostname=0.0.0.0" ]
