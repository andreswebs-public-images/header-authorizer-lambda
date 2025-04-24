package main

import (
	"context"
	"net/http"
	"os"

	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var (
	cfg       *aws.Config
	ssmClient *ssm.Client
	env       slog.Attr

	principalIDEnvVar      string = "PRINCIPAL_ID"
	headerKeyEnvVar        string = "HEADER_KEY"
	headerValueParamEnvVar string = "HEADER_VALUE_PARAMETER"

	principalID          string
	headerKey            string
	headerValueParamName string
	secret               string
)

type Request events.APIGatewayCustomAuthorizerRequestTypeRequest
type Response events.APIGatewayCustomAuthorizerResponse

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	principalID = ReadEnvVarWithDefault(principalIDEnvVar, "user")
	headerKey = ReadRequiredEnvVar(headerKeyEnvVar)
	headerValueParamName = ReadRequiredEnvVar(headerValueParamEnvVar)

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.Error("unable to load AWS configuration", slog.Any("err", err))
		os.Exit(1)
	}

	ssmClient = ssm.NewFromConfig(cfg)

	headerValueParam, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(headerValueParamName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		slog.Error("failed to fetch SSM parameter", slog.String("parameterName", headerValueParamName), slog.Any("err", err))
		os.Exit(1)
	}

	secret = *headerValueParam.Parameter.Value

	env = slog.Group(
		"env",
		slog.String(principalIDEnvVar, principalID),
		slog.String(headerKeyEnvVar, headerKey),
		slog.String(headerValueParamEnvVar, headerValueParamName),
	)
}

func handler(ctx context.Context, req Request) (Response, error) {

	// handle case sensitive headers
	// https://github.com/aws/aws-lambda-go/issues/117
	headers := http.Header{}
	for header, value := range req.Headers {
		headers.Add(header, value)
	}

	headerValue := headers.Get(headerKey)
	ok := headerValue != ""

	if headerValue == secret {
		slog.Info("allowed", slog.String("type", "application"), env, slog.Bool("h", ok), slog.String("target", req.MethodArn))
		return generatePolicy(principalID, "Allow", req.MethodArn), nil
	}

	slog.Info("denied", slog.String("type", "application"), env, slog.Bool("h", ok), slog.String("target", req.MethodArn))
	return generatePolicy(principalID, "Deny", req.MethodArn), nil
}

func generatePolicy(principalID, effect, resource string) Response {
	return Response{
		PrincipalID: principalID,
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Effect:   effect,
					Action:   []string{"execute-api:Invoke"},
					Resource: []string{resource},
				},
			},
		},
	}
}

func main() {
	lambda.Start(handler)
}
