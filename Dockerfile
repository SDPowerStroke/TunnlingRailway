# Base image with Node
FROM node:20-bullseye

# Install Python + pip
RUN apt-get update && \
    apt-get install -y python3 python3-pip && \
    rm -rf /var/lib/apt/lists/*

# Install yt-dlp
RUN pip3 install --no-cache-dir yt-dlp

# App directory
WORKDIR /app

# Copy files
COPY package.json ./
RUN npm install

COPY . .

# Railway sets PORT automatically
EXPOSE 3000

CMD ["npm", "start"]
