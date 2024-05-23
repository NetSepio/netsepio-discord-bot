# Use the official Node.js image.
# The version of Node.js can be adjusted as necessary.
FROM node:14

# Create and change to the app directory.
WORKDIR /usr/src/app

# Copy package.json and package-lock.json.
COPY package*.json ./

# Install dependencies.
RUN npm install

# Copy the rest of the application code.
COPY . .

# Expose the port the app runs on, if needed (for example, 8080).
# EXPOSE 8080

# Define environment variable for the path to .env file.
# ENV NODE_ENV=production

# Command to run the application.
CMD ["node", "index.js"]

# Note: Replace 'index.js' with the actual entry point of your application.
