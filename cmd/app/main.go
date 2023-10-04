package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ServerlessChatgptWrapperStackProps struct {
	awscdk.StackProps
}

func NewServerlessChatgptWrapperStack(scope constructs.Construct, id string, props *ServerlessChatgptWrapperStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	api := awsapigateway.NewRestApi(stack, jsii.String("ChatGPTApiWrapper"), &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins:     awsapigateway.Cors_ALL_ORIGINS(),
			AllowHeaders:     awsapigateway.Cors_DEFAULT_HEADERS(),
			AllowMethods:     awsapigateway.Cors_ALL_METHODS(),
			AllowCredentials: jsii.Bool(true),
		},
		CloudWatchRole: jsii.Bool(true),
	})

	testLambda := awslambda.NewFunction(stack, jsii.String("TestLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("./lambdas/test"), &awss3assets.AssetOptions{}),
	})

	testResource := api.Root().AddResource(jsii.String("test"), &awsapigateway.ResourceOptions{})

	testIntegration := awsapigateway.NewLambdaIntegration(testLambda, &awsapigateway.LambdaIntegrationOptions{})

	testResource.AddMethod(jsii.String("POST"), testIntegration, &awsapigateway.MethodOptions{
		ApiKeyRequired: jsii.Bool(true),
	})

	testKey := api.AddApiKey(jsii.String("testApiKey"), nil)

	testLambda.AddPermission(jsii.String("AllowAPIGatewayInvoke"), &awslambda.Permission{
		Action:    jsii.String("lambda:InvokeFunction"),
		Principal: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		SourceArn: api.ArnForExecuteApi(jsii.String("*"), jsii.String("/test"), api.DeploymentStage().StageArn()),
	})

	usagePlan := awsapigateway.NewUsagePlan(stack, jsii.String("testUsagePlan"), &awsapigateway.UsagePlanProps{
		Name: jsii.String("testUsagePlan"),
		Throttle: &awsapigateway.ThrottleSettings{
			RateLimit:  jsii.Number(10),
			BurstLimit: jsii.Number(2),
		},
	})

	usagePlan.AddApiKey(testKey, nil)

	usagePlan.AddApiStage(&awsapigateway.UsagePlanPerApiStage{
		Api:   api,
		Stage: api.DeploymentStage(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewServerlessChatgptWrapperStack(app, "ServerlessChatgptWrapperStack", &ServerlessChatgptWrapperStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
