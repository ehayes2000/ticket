FROM node:latest as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY ./ .
RUN npm run build

FROM nginx as production-stage
RUN mkdir /service
COPY --from=build-stage /app/dist /service
COPY nginx.conf /etc/nginx/nginx.conf


