# AWS Bedrock Chat CLI

This is a small command-line utility written in Go to make API requests to the AWS Bedrock chat API for Claude.
The tool takes an optional model as input and gives a `Charm.sh`'s `huh` based interface which lets you query Claude.
On exit it will show tokens consumed and timing metrics.

## Usage

To use the CLI tool, you need to have the your AWS environment variables setup as usual (e.g. a working AWS_PROFILE) that allows for `bedrock:InvokeModel` permission.
Also you'll need to ask AWS to allow you to use these models, basically just go to [AWS Bedrock](https://us-east-1.console.aws.amazon.com/bedrock/home?region=us-east-1#/modelaccess) and request access to all their models.
It's free (getting access) and usually doesn't take very long (minutes/hours rather than days).

```
$ export AWS_PROFILE=bedrock
$ ./chat
```

The `--model` flag can be used to specify a different model (default is `anthropic.claude-v2`, try `anthropic.claude-3-sonnet-20240229-v1:0` for more expensive chat, see ):

```
$ ./chat --model "anthropic.claude-3-sonnet-20240229-v1:0"
```

## Building

To build the CLI tool, you can either use your local golang installation (needs 1.20 or higher), or have Docker build it.

### Docker

Run the following command to build the Docker image:

```
$ docker build -t chat:latest .
```

After the build is complete, you can run the CLI tool using the Docker image:

```
$ docker run --rm -e AWS_PROFILE=<yourprofile> -v ~/.aws:/root/.aws:ro chat:latest
```
Note that the AWS credentials need to be present in some form, if you have setup a local profile this will do. See https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html for details.
E.g. alternatively you can provide the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` and `AWS_DEFAULT_REGION` variables or simply run it from e.g. EC2 with a proper instance profile.

Anyway, above can also be done using make:
```
make build-docker docker-run
```

### Local build

Run the following command to build the chat binary:
```
$ make build-local
```

After building you can run it directly:
```
$ ./chat
```

## Notes

This is just a pet project to tinker with [Charm](https://charm.sh/) and Claude.

