# build environment
FROM node:12.2.0-alpine as build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY app/package.json /app/package.json
RUN npm install --silent
RUN npm install react-scripts@3.0.1 -g --silent
COPY app /app/

ENV CONVERTER_SERVER converterserver:3001

RUN npm run build

FROM mhart/alpine-node
RUN yarn global add serve
RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build /app/build .
EXPOSE 80
CMD ["serve", "-p", "80", "-s", "."]


## production environment
#FROM nginx:1.16.0-alpine
#COPY --from=build /app/build /usr/share/nginx/html
## Copy custom configuration file from the current directory
#COPY nginx.conf /etc/nginx/nginx.conf
#
#EXPOSE 80
#
#ENV CONVERTER_SERVER converterserver:3001
#CMD ["nginx", "-g", "daemon off;"]
#