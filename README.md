# docker-plugin-s3

Simple Docker CLI plugin that supports reading/writing to S3 compatible storage.

This is **not** meant to be a replacement for OCI-compatible registries but to provide a simple CLI API to move locally stored images to and from blob storage.

If you at all need to think about access management or container lifecycle management, invest in buying or hosting your own OCI registry!

## Usage

### Dependencies

- Golang `v1.20.2`
- Docker Client `v1.41`

### Installation

From within the repo:

- Run `make`.
- Run `docker help`, you should see `s3*` under the section `Management Commands`.

### Commands

- `docker s3 push`: Push a Docker image to S3 compatible storage
- `docker s3 pull`: Pull a Docker image from S3 compatible storage
- `docker s3 docker-cli-plugin-metadata`: Exposes plugin metadata needed to register the CLI plugin with the Docker CLI.

### Non-standard S3 Endpoint

If you need to change the default AWS-provided S3 endpoint, set the variable `AWS_S3_ENDPOINT` in your environment.

### Local Development

This repo provides a simple docker-compose file to run [MinIO](https://github.com/minio/minio) (an open-source S3 compatible storage solution) locally.
It can be stood up using the command `docker compose -f dev/s3/docker-compose.yaml up`.

This repo also provides a `.local.env` file. You can use a tool like [dotenv-cli](https://www.npmjs.com/package/dotenv-cli) like so `dotenv -e .local.env -- docker s3 $SOME_SUBCOMMAND` to export it temporarily into your command's environment.
