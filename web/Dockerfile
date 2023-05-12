# Build weave frontend
FROM node:lts-alpine AS builder

ARG BUILD_OPTS=build

WORKDIR /weave
COPY . .
RUN npm install && \
    npm run ${BUILD_OPTS}

# Setup nginx server for frontend
FROM nginx:stable-alpine
COPY --from=builder /weave/dist/ /usr/share/nginx/html/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;" ]
