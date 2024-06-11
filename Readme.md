# Redmage

Redmage is a Reddit Image Downloader with Scheduler, Image Filters, and GUI for management.
Redmage is self-hosted friendly with single binary and docker image.

## Features

1. Image Filter by Resolution and Aspect Ratio.
2. Device profiles for different devices.
   Each device will have its own list of images and filters.
3. Scheduler to download images from Subreddits at specific time.
4. Web GUI for management.
5. Images are stored in normal directories and can be accessed easily. Syncs well with cloud storage or Syncthing.

## Installation

### Binary

Download from [Releases] page, select your OS and architecture, and run the binary.

### Docker

```bash
docker run  -d -p 8080:8080 -v /path/to/images:/app/images -v /path/to/database:/app/db git.tigor.web.id/tigor/redmage
```

### Docker Compose

See [docker-compose.yml](./docker-compose.yaml) for example.

## Usage

After installation, go to `localhost:8080` to access the web GUI. Add a Device Profile and subscribe to Subreddits. You can also set the Scheduler to download images at specific time.
