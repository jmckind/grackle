# grackle

A simple framework for ingesting and displaying Tweets.

## Overview

Grackle was primarily designed as a demo [Go](https://golang.org/) application for making use of the
[RethinkDB](https://www.rethinkdb.com/) database. The application offers two binary components `grackle-ingest`
and `grackle-web`.

### grackle-ingest

The `grackle-ingest` component will stream Tweets using the Twitter API, based on one or more [search terms][twitter_api_track],
and store this data in a RethinkDB database.

### grackle-web

The `grackle-web` component runs a basic web UI for displaying the incoming Tweet data from the RethinkDB database in
real-time using a [Changefeed][rethink_changefeed].

## Usage

Grackle was primarily intended to be deployed on Kubernetes and manifests for this are located in the `deploy` directory.

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

An example RethinkDB deployment manifest for testing and local development is located in the `deploy` directory.

```bash
kubectl create -f deploy/rethinkdb.yaml
```

### Twitter Developer Credentials

Twitter API credentials will be needed before deployment. Either create a new Twitter
[developer application](https://developer.twitter.com/en/apps) or use an existing appilcation and navigate to the
`Keys and tokens` tab to grab the credentials.

### Secret

Update the Secret manifest `deploy/secret.yaml` with the credentials and remember that the values must be base64 encoded.

```bash
echo -n "my-twitter-access-token" | base64 -
```

Create the secret resource in the new namespace.

```bash
kubectl create -f deploy/secret.yaml
```

### Ingest

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

A working Go environment and GNU Make are required. There are various Make targets for building the application.

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
[twitter_api_track]:https://developer.twitter.com/en/docs/tweets/filter-realtime/guides/basic-stream-parameters#track
[rethink_changefeed]:https://www.rethinkdb.com/docs/changefeeds/javascript/
