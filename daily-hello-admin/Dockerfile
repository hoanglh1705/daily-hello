FROM node:20-alpine AS build

WORKDIR /app

# Copy package configurations
COPY package.json package-lock.json* ./

# Install dependencies cleanly
RUN npm install

# Copy source code and build
COPY . .
RUN npm run build

FROM nginx:alpine

# Copy built app to nginx html folder
COPY --from=build /app/dist /usr/share/nginx/html

# Add a default mapping configuration to support React Router SPA Fallback
RUN rm /etc/nginx/conf.d/default.conf
RUN echo "server { \
    listen 80; \
    server_name localhost; \
    location / { \
        root /usr/share/nginx/html; \
        index index.html index.htm; \
        try_files \$uri \$uri/ /index.html; \
    } \
}" > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
