[![Go Report Card](https://goreportcard.com/badge/github.com/theopenlane/openlane-cloud)](https://goreportcard.com/report/github.com/theopenlane/openlane-cloud)
[![Build status](https://badge.buildkite.com/9d99bb1f92d9195776d9983bea1f74314fd912706244c48863.svg)](https://buildkite.com/theopenlane/theopenlane-cloud)
[![Go Reference](https://pkg.go.dev/badge/github.com/theopenlane/openlane-cloud.svg)](https://pkg.go.dev/github.com/theopenlane/openlane-cloud)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache2.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)

# Openlane Cloud

Building a SaaS offering on top of [Openlane](https://github.com/theopenlane/core)

## Openlane Cloud Server

The Openlane Cloud server is used to consume the [openlane server](https://github.com/theopenlane/core/) and apply an opinionated implementation on top of the the generics provided. Many, if not all, of the endpoints provided by the server use the [openlane client](https://github.com/theopenlane/core/blob/main/pkg/openlaneclient/client.go) to make requests to the openlane Server.

As an example, the `v1/organizations` endpoint uses the `openlane` client to create an organizational hierarchy`:

```
│   └── rootorg <--- top level organization
│       ├── production <-- top level environment per customer organization
│       │   ├── assets <-- buckets
│       │   ├── customers
│       │   ├── orders
│       │   ├── relationships
│       │   │   ├── internal_users  <-- relationships
│       │   │   ├── marketing_subscribers
│       │   │   ├── marketplaces
│       │   │   ├── partners
│       │   │   └── vendors
│       │   └── sales
│       └── test <-- organization identical to production just named
```

## Openlane Cloud CLI

The openlane cloud cli is used to interact with the openlane cloud server as well as some requests directly to the openlane server using the [openlane client](https://github.com/theopenlane/core/blob/main/pkg/openlaneclient/client.go). In order to use the cli, you must have a registered user with the openlane server.

### Installation

```bash
brew install theopenlane/tap/openlane-cloud
```

### Upgrade

```bash
brew upgrade theopenlane/tap/openlane-cloud
```

### Usage

```bash
openlane-cloud
the openlane-cloud cli

Usage:
  openlane-cloud [command]

Available Commands:
  completion     Generate the autocompletion script for the specified shell
  help           Help about any command
  seed           the subcommands for creating demo data in openlane
  organization   the subcommands for working with the openlane organization
```

## Seeding Data

The `openlane-cloud` cli has functionality to generate and load test data into `openlane` using the `seed` command.

```bash
Usage:
  openlane-cloud seed [command]

Available Commands:
  generate    generate random data for seeded environment
  init        init a new openlane seeded environment
```

#### Using the Taskfile

On a brand new database, you should run:

1. Create a new user to authenticate with the openlane API, this command will fail on subsequent tries because the user will already exist.
    ```bash
    task register
    ```
1. Login as the user, create a new Personal Access Token that will be used to seed the data, generate a new data set, bulk load objects into the openlane API:
    ```bash
    task cli:seed:all
    ```

If instead, you prefer to use the CLI commands directly, keep reading.

### Generate Data

Using the `generate` subcommand, new random data will be stored in csv files:

```bash
openlane-cloud seed generate
```

<details>
<summary>Generated Data</summary>

```bash
tree demodata
demodata
├── groups.csv
├── invites.csv
├── orgs.csv
└── users.csv
```

</details>

### Init Environment

Using the `init` subcommand, the data in the specified directory (defaults to `demodata` in the current directory), the csv files will be used to generate the data.

```bash
openlane-cloud seed init
```

The newly created objects will be displayed when complete:

<details>
<summary>Results</summary>

```bash
> seeded environment created 100% [===============]  [3s]
Seeded Environment Created:
+--------------------------------------------------------------------------------------+
| Organization                                                                         |
+----------------------------+--------+-------------+-------------+----------+---------+
| ID                         | NAME   | DESCRIPTION | PERSONALORG | CHILDREN | MEMBERS |
+----------------------------+--------+-------------+-------------+----------+---------+
| 01J06RPZ8HQRWW4AZERHKWT2YH | Plus-U |             | false       |        0 |       1 |
+----------------------------+--------+-------------+-------------+----------+---------+
...
```

</details>

## Contributing

Please read the [contributing](.github/CONTRIBUTING.md) guide as well as the [Developer Certificate of Origin](https://developercertificate.org/). You will be required to sign all commits to the Openlane project, so if you're unfamiliar with how to set that up, see [github's documentation](https://docs.github.com/en/authentication/managing-commit-signature-verification/about-commit-signature-verification).

## Licensing

This repository contains `openlane-cloud` which is open source software under [Apache 2.0](LICENSE). Openlane is a product produced from this open source software exclusively by theopenlane, Inc. This product is produced under our published commercial terms (which are subject to change), and any logos or trademarks in this repository or the broader [theopenlane](https://github.com/theopenlane) organization are not covered under the Apache License.

Others are allowed to make their own distribution of this software or include this software in other commercial offerings, but cannot use any of the Openlane logos, trademarks, cloud services, etc.

## Security

We take the security of our software products and services seriously, including all of the open source code repositories managed through our Github Organizations, such as [theopenlane](https://github.com/theopenlane). If you believe you have found a security vulnerability in any of our repositories, please report it to us through coordinated disclosure.

**Please do NOT report security vulnerabilities through public github issues, discussions, or pull requests!**

Instead, please send an email to `security@openlane.io` with as much information as possible to best help us understand and resolve the issues. See the security policy attached to this repository for more details.

## Questions?

You can email us at `info@openlane.io`, open a github issue in this repository, or reach out to [matoszz](https://github.com/matoszz) directly.



