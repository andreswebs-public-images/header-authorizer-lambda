# header-authorizer

AWS Lambda Authorizer for API Gateway that checks for a secret in a custom
header.

## Configuration and Pre-requisites

This lambda requires a few environment variables:

- `HEADER_KEY`: The value of this variable is the name of the HTTP header that
  the lambda will check.
- `HEADER_VALUE_PARAMETER`: The value of this variable is the name of an SSM
  Parameter containing the secret that the lambda will look for in the specified
  header.

It also accepts an optional environment variable:

- `PRINCIPAL_ID`: The value to use for the `principalId` field in the lambda
  output. Defaults to `user`.

The SSM parameter specified in `HEADER_VALUE_PARAMETER` must exist with the
correct value for the lambda to work.

The IAM role assigned to the lambda must have permissions to get and decrypt
this SSM parameter.

## Authors

**Andre Silva** - [@andreswebs](https://github.com/andreswebs)

## License

This project is licensed under the [Unlicense](UNLICENSE.md).

## References

<https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-use-lambda-authorizer.html>

<https://aws.amazon.com/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/>
