# grackle

A simple cloud-native framework for ingesting and displaying tweets.

## Usage

Kubernetes manifests are located in the `deploy` directory.

```bash
$ tree deploy
deploy
├── ingest.yaml
├── rethinkdb.yaml
├── secret.yaml
└── web.yaml
```

### Namespace

Create a new namespace for the application.

```bash
kubectl create namespace grackle
```

### RethinkDB Cluster

An example deployment for RethinkDB located in the `deploy` directory can be used for testing.

```bash
kubectl create -f deploy/rethinkdb.yaml
```

### Twitter Developer Credentials

Twitter API credentials will be needed before deployment. Either create a new
Twitter developer application or use an existing appilcation and navigate to the `Keys and tokens` tab to grab the
credentials.

<https://developer.twitter.com/en/apps>

### Secret

Update the Secret manifest `deploy/secret.yaml` with the credentials and remember that the value must be base64 encoded.

```bash
echo "my-twitter-access-token" | base64 -
```

Create the secret resource in the new namespace.

```bash
kubectl create -f deploy/secret.yaml
```

### Ingest Process

Once the RethinkDB database is running and the Secret reseource has been created, create the deployment for the Twitter
ingest process.

```bash
kubectl create -f deploy/ingest.yaml
```

### Web UI

Once the ingest process is running, create the deployment for the web UI.

```bash
kubectl create -f deploy/web.yaml
```

## Development

A working golang environment and GNU Make are required. There are various Make target for building the application.

Build the application.

```bash
make build
```

Run the ingest process locally.

```bash
make run-local-grackle-ingest
```

## License

Grackle is released under the Apache 2.0 license. See the [LICENSE][license_file] file for details.


[license_file]:./LICENSE
