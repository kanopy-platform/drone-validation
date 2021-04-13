# Validation extension for Drone

_This project requires Drone server version 1.4 or higher._

## Installation

Create a shared secret:

```shell
openssl rand -hex 16
```

Download and run the plugin:

```shell
docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=<your_shared_secret> \
  --restart=always \
  --name=<container_name> quay.io/mongodb-labs/drone-validation:latest 
```

Update your `drone-server`` environment variables to include the plugin endpoint and shared secret.

```shell
DRONE_VALIDATE_PLUGIN_ENDPOINT=http://<service-name>:3000
DRONE_VALIDATE_PLUGIN_SECRET=<your_shared_secret>
```
