# secretsauce

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/ajmilazzo/secretsauce)

![spicy](https://github.com/ajmilazzo/secretsauce/assets/18070948/56f818ce-7b65-4a5f-8128-c46ecf8566ae)

## Overview

SecretSauce is a Kubernetes operator built to manage and govern the creation and lifecycle of randomly generated secrets within a Kubernetes cluster. Developed with `kubebuilder`, the project showcases a Kubernetes controller, admission webhook, and custom resource definitions (CRDs). The operator is focused on security and compliance, enabling both the dynamic creation of random secrets and setting governance rules for these secrets.

## Features:

- **RandomSecret CRD**: Define and create random Kubernetes secrets with configurable parameters like name and length.
- **SecretPolicy CRD**: Set policies to enforce certain requirements, such as minimum length, on RandomSecrets.
- **RandomSecret Controller**: Handles the creation and lifecycle management of RandomSecrets.
- **SecretPolicy Admission Webhook**: Ensures secrets adhere to the established SecretPolicy before their creation.

**RandomSecret CRD**
```
apiVersion: fancysecrets.secretsauce.anthonymilazzo.com/v1
kind: RandomSecret
metadata:
  name: randomsecret-sample
spec:
  secretName: test-rand
  length: 16
```

**SecretPolicy CRD**
```
apiVersion: fancysecrets.secretsauce.anthonymilazzo.com/v1
kind: SecretPolicy
metadata:
  name: secretpolicy-sample
spec:
  minLength: 12
```

## Requirements:

- Kubernetes cluster
- cert-manager (for facilitating TLS for the webhook)

## Design Decisions & Rationale:

- **Use of `kubebuilder`**: Chosen for its extensive documentation, rich feature set, and popularity within the Kubernetes community. While frameworks might abstract some of the intricacies, it was chosen to expedite development within the provided timeline.
- **Security**: Emphasis on a cryptographically secure method to generate random secrets. This ensures secrets are robust and unpredictable.
- **SecretPolicy Admission Webhook**: Provides an extra layer of governance to ensure compliance and consistency across secrets.

## Business Logic:
- **`internal/controller`**: Contains the logic for creating and managing `RandomSecret` resources within the cluster
- **`internal/webhook`**: Contains the logic for the `SecretPolicy` admission webhook, which allows users to govern the minimum length of the `RandomSecret` resources
- **`cmd/main.go`**: The entrypoint, which configures the webhook and controller

## Future Enhancements:

- **Complexity Setting**: Allow users to define the complexity of secrets, e.g., mix of characters, numbers, symbols, etc.
- **Minimum Complexity Requirements**: Ensure that secrets match a certain complexity pattern or standard.
- **Secret Rotation**: Choose how often secrets are rotated.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster. A `k3d` cluster will be automatically provisioned if you open this repo in `Codespaces` or `devcontainer`.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

## One-step setup (k3d)
If you are using k3d, and do not yet have a cluster setup, the following command will bring up the cluster, cert-manager, registry, and install the controller

```sh
make k3d-all
```

If you are using the Codespace environment (which already has a cluster setup), use the following command to build, push, and deploy the controller

```sh
make k3d-build-deploy
```

## Setting up cert-manager
Before deployment, ensure you have `cert-manager` running. It can be installed with:

```sh
make certmanager
```

## Setting up the cluster (k3d)
**Note:** The following steps assume you are using `k3d`, other clusters will differ. These steps are automatically performed if using the `Codespaces` environment.

Setup the k3d cluster and image registry:

```sh
make create
```

If you are getting the error `resolving host k3d-registry.localhost: lookup k3d-registry.localhost: no such host` then you may need to add the following entry to `/etc/hosts`:

```sh
127.0.0.1 k3d-registry.localhost
```

## Running on a generic cluster
1. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/secretsauce:tag
```

2. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/secretsauce:tag
```

## Running on a k3d cluster
1. Build and push your image to the location specified by `IMG`:

```sh
make docker-build-k3d docker-push-k3d
```

2. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy-k3d
```
## Deploying sample CRDs

Install Instances of Custom Resources:

**SecretPolicy CRD**
```sh
kubectl apply -k config/samples/policy/
```

**RandomSecret CRD**
```sh
kubectl apply -k config/samples/secret/
```

You should see a success and error (for the example that doesn't meet the SecretPolicy):

```
randomsecret.fancysecrets.secretsauce.anthonymilazzo.com/randomsecret-sample-good created
Error from server (Forbidden): error when creating "config/samples/secret/": admission webhook "vrandomsecret.kb.io" denied the request: RandomSecret length 10 is less than SecretPolicy MinLength 12
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

### Delete the cluster (k3d)
Delete the cluster and registry

```sh
make delete
```

## How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

## Run Locally
1. Install the CRDs into the cluster:

```sh
make install
```

2. Ensure TLS certs are available locally running:

```sh
make tls
```

3. Run the controller locally:

```sh
make run
```

## Running Tests
Tests can be run with:

```sh
make test
```

## Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## Attribution
Credit to the following for providing templates and scaffolding for this repo:
- [kubernetes-sigs/kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) - Scaffolding of the controller and manifests
- [cse-labs/kubernetetes-in-codespaces](https://github.com/cse-labs/kubernetes-in-codespaces) - Codespaces template for k3d
- [cse-labs/codespaces-images](https://github.com/cse-labs/codespaces-images) - Codespaces docker images

## License

Copyright 2023 Anthony Milazzo.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

