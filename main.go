package main

import (
	"context"
	"net/http"
	"os"

	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	sess                 *session.Session
	ssmClient            *ssm.SSM
	principalID          string
	headerKey            string
	headerValueParamName string
	secret               string
	env                  slog.Attr
)

type Request events.APIGatewayCustomAuthorizerRequestTypeRequest
type Response events.APIGatewayCustomAuthorizerResponse

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	sess = session.Must(session.NewSession())
	ssmClient = ssm.New(sess)

	principalID = os.Getenv("PRINCIPAL_ID")
	if principalID == "" {
		principalID = "user"
	}

	headerKey = os.Getenv("HEADER_KEY")
	if headerKey == "" {
		slog.Error("missing environment variable HEADER_KEY")
		os.Exit(1)
	}

	headerValueParamName = os.Getenv("HEADER_VALUE_PARAMETER")
	if headerValueParamName == "" {
		slog.Error("missing environment variable HEADER_VALUE_PARAMETER")
		os.Exit(1)
	}

	headerValueParam, err := ssmClient.GetParameter(&ssm.GetParameterInput{
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
		slog.String("PRINCIPAL_ID", principalID),
		slog.String("HEADER_KEY", headerKey),
		slog.String("HEADER_VALUE_PARAMETER", headerValueParamName),
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
