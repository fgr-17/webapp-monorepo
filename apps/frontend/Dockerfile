# Quality gate (ESLint + Prettier) then ship static assets via nginx.
FROM node:20-alpine AS quality

WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY eslint.config.js .prettierrc.json .prettierignore ./
COPY index.html app.js calculator.js calculator.test.js styles.css ./
RUN npm run lint

FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=quality /src/index.html /src/app.js /src/calculator.js /src/styles.css /usr/share/nginx/html/

EXPOSE 80
