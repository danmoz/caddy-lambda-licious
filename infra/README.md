# AWS Test Infrastructure

Sets up real AWS infrastructure to test this caddy plugin.

Running `mise aws-deploy` will:

- Build the python FastAPI test app into a lambda.zip deployment package
- Create a ECR docker repo service
- Build the container image locally with `xcaddy`
- Push the image to the ECR repo
- Set up an ECS cluster the deploys the image from ECR
- Sets up a lambda running the python service

The service is intentionally HTTP-only for testing; so no SSL.

The stack creates a VPC named `caddy-lambda-licious-testing-vpc` with two public
subnets and an internet gateway. It does not require a default VPC.

## Prerequisites

- [mise](https://mise.jdx.dev/)
- Docker running locally
- AWS credentials in your shell (run `aws login`)

## Deploy

Trust the nested mise configuration once:

```sh
mise trust infra/mise.toml
```

Then run the commands from this directory:

```sh
cd infra
mise aws-deploy
```

The default region is `us-east-1`, configured as `AWS_REGION` in
`infra/mise.toml`. Override it and resource names with environment variables:

```sh
AWS_REGION=ap-southeast-2 mise aws-deploy
```

The stack defaults to `api_gateway_v1`. Caddy receives the Terraform
`event_format` variable through the ECS task environment. To compare against v2:

```sh
TF_VAR_event_format=api_gateway_v2 mise aws-deploy
```

Run `mise aws-deploy` again to rebuild the image from your current checkout and
deploy changes to an existing stack. Deployment prints the public Caddy URL.
Use that URL to test repeated query parameters, cookies, and binary responses:

```sh
curl 'http://PUBLIC_IP/echo?tag=one&tag=two' \
  -H 'Cookie: session=abc' -H 'X-Test: one' -H 'X-Test: two'
curl -i http://PUBLIC_IP/cookies
curl -s http://PUBLIC_IP/binary | od -An -tx1
```

The echo response includes the raw event headers and `multiValueHeaders` so you
can inspect the v1 representation. The binary response should contain
`00 01 fe ff`.

## Destroy

Remove the test resources, including the ECR repository and image:

```sh
mise aws-destroy
```
