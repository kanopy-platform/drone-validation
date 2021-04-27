# Validation extension for Drone

_Requires Drone server version 1.4 or higher._

This is a simple validation extension that performs pipeline type validations.

When a user prepared a Drone configuration file using a pipeline that was not supported by our platform, the build used to hang and users did not receive any error or feedback. Whith this validation extension enabled, the same build would instantly fail and an error message returned to the UI informing the reason why the build has failed.

## Implementation details

The validation process relies on an [OPA policy](https://www.openpolicyagent.org/) written in [rego](https://www.openpolicyagent.org/docs/latest/policy-language/) and is integrated via [OPA Go API](https://www.openpolicyagent.org/docs/latest/integration/#integrating-with-the-go-api).

Although the main purpose of this extension's policy is to validate supported pipeline types, it can be easily modified to check different attributes of the Drone configuration, including secrets.

This application runs a simple rego query that requires the policy to return a boolean value on the path `data.drone.validation.deny` and an error message on `data.drone.validation.out`.

The [default policy](policy/validation.rego) is based on the following workflow:

``` text

                 ┌───────────────────┐         ┌──────────────────────┐        ┌───────────────────┐
 Config document │                   │  Yes    │                      │  No    │                   │
      ──────────►│ Is it a pipeline? ├────────►│ Is the type allowed? ├───────►│ Return: deny=true │
                 │                   │         │                      │        │                   │
                 └─────────┬─────────┘         └───────────┬──────────┘        └───────────────────┘
                           │                               │
                           │ No                            │ Yes
                           │                               │
                           ▼                               ▼
                ┌────────────────────┐          ┌────────────────────┐
                │                    │          │                    │
                │ Return: deny=false │          │ Return: deny=false │
                │                    │          │                    │
                └────────────────────┘          └────────────────────┘

```

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

## Auditing

Drone doesn't offer audit logging at the moment, so we can use this validation extension to generate log entries for both build and promotion job executions.

Log messages are opinionated and don't contain the full list of attributes by default.

