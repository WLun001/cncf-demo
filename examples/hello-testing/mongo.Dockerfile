FROM mongo:5
COPY mongo/init.js /docker-entrypoint-initdb.d/init.js
