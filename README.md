# docker-plugin-s3

Simple docker plugin that supports reading/writing to S3 compatible storage.`

## Installation

Clone the repo, run `make`!

## Usage

### Non-standard S3 Endpoint

If you need to change the S3 Endpoint (for ex.: you use Minio) set the variable `AWS_S3_ENDPOINT` in your environment.

### Commands

- `docker s3 push`: Push a Docker image to S3 compatible storage
- `docker s3 pull`: Pull a Docker image from S3 compatible storage
