A Drone validation extension for the Kanopy platform.

_This project requires Drone server version 1.4 or higher._

## Installation

Create a shared secret:

```shell
openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```shell
docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --restart=always \
  --name=moo registry.example.com/drone-validation
```

If you are deploying locally using helm:

```shell
helm install -n drone-validation -f ./environment/dev.yaml mongodb/web-app
```

Update your Drone server environment variables to include the plugin address and the shared secret.

```shell
DRONE_VALIDATE_PLUGIN_ENDPOINT=http://<service-name>:3000
DRONE_VALIDATE_PLUGIN_SECRET=bea26a2221fd8090ea38720fc445eca6
```

## Kanopy Deployment

Access the dronev1 namespace in the Kanopy cluster where you want to enable the validation extension.

Create a new value in drone-secrets with the validate plugin secret:

``` shell
ksec set drone-secrets DRONE_VALIDATE_PLUGIN_SECRET=a750ac1f1a411197d7ae930cd5d01fc
```
