# Husky Watch

## Overview

This website was developed to help University of Washington Tacoma students be informed about the crime around campus. This can easily be expanded to include Seattle, Bothell (if they make their crimes easily viewable and exportable), and any other city in the country!

The client directory contains the frontend. The server directory contains the backend.
The data_parser directory contains a small go program that transferred the tacoma csv into our database.
The sql directory contains the the crime data in csv as well as some sql scripts.

## Technologies

- Svelte + tailwindcss
- Go + Gin framework
- PostgreSQL
- Docker

## How to use

***A GOOGLE MAPS API KEY IS REQUIRED FOR THE MAP PAGE TO WORK!!!**

If using the the env below, the server should be able to connect to the database.

To run with docker compose, a .env file is required in the project's root directory containing:

```
VITE_PORT=5173
PUBLIC_API_URL=http://server:4000
PUBLIC_GOOGLE_MAPS_API_KEY=
PUBLIC_GOOGLE_MAPS_ID=
SERVER_PORT=4000

NODE_ENV=production
GIN_MODE=release

POSTGRES_URL=
```

