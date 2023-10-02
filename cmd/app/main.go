package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
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
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
			AllowHeaders: jsii.Strings(*jsii.String("X-Api-Key")),
		},
	})

	testLambda := awslambda.NewFunction(stack, jsii.String("TestLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("./lambdas/test"), nil),
	})
	testResource := api.Root().AddResource(jsii.String("test"), nil)
	testIntegration := awsapigateway.NewLambdaIntegration(testLambda, nil)
	testResource.AddMethod(jsii.String("POST"), testIntegration, &awsapigateway.MethodOptions{
		ApiKeyRequired:    jsii.Bool(true),
		AuthorizationType: awsapigateway.AuthorizationType_NONE,
	})
	testKey := api.AddApiKey(jsii.String("testApiKey"), nil)
	usagePlan := awsapigateway.NewUsagePlan(stack, jsii.String("testUsagePlan"), &awsapigateway.UsagePlanProps{
		Name: jsii.String("testUsagePlan"),
		Throttle: &awsapigateway.ThrottleSettings{
			RateLimit:  jsii.Number(10),
			BurstLimit: jsii.Number(2),
		},
	})
	usagePlan.AddApiKey(testKey, nil)

	deployment := awsapigateway.NewDeployment(stack, jsii.String("ApiDeployment"), &awsapigateway.DeploymentProps{
		Api: api,
	})

	stage := awsapigateway.NewStage(stack, jsii.String("ProdStageV2"), &awsapigateway.StageProps{
		Deployment: deployment,
		StageName:  jsii.String("prodV2"),
	})

	usagePlan.AddApiStage(&awsapigateway.UsagePlanPerApiStage{
		Api:   api,
		Stage: stage,
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
