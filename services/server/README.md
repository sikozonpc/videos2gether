# Videos Together server

Server source code for the https://videos2gether.com

### Architecture

The server aims to be as slim as possible since it's deployed as a cloud run service that is booted up on room creation.

In order to provide the super quick boot time it's deployed a binary file built from the Golang compiler.

The server acts as a WebSocket service for the realtime actions for the rooms and an HTTP service for serving some endpoints (although in the future I aim to remove the HTTP layer in favor of all WebSocket).

### CI/CD

When a commit is merged on master a GCloud trigger is set to run the `cloudbuild.yaml` file and deploy a new revision of the cloud run service.