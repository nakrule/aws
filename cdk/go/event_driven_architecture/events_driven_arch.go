package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type EventsDrivenArchStackProps struct {
	awscdk.StackProps
}

func NewEventsDrivenArchStack(scope constructs.Construct, id string, props *EventsDrivenArchStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	// Lambda functions processing event from an EventBridge Bus.
	fn_a := awslambda.NewFunction(stack, jsii.String("Function A"), &awslambda.FunctionProps{
		FunctionName: jsii.String("Function-A"),
		Runtime:      awslambda.Runtime_GO_1_X(),
		Handler:      jsii.String("main"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/function_a/"), nil),
	})

	fn_b := awslambda.NewFunction(stack, jsii.String("Function B"), &awslambda.FunctionProps{
		FunctionName: jsii.String("Function-B"),
		Runtime:      awslambda.Runtime_GO_1_X(),
		Handler:      jsii.String("main"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/function_b/"), nil),
	})

	// EventBridge bus
	bus := awsevents.NewEventBus(stack, jsii.String("bus"), &awsevents.EventBusProps{
		EventBusName: jsii.String("cdk-stack-bus-lambda"),
	})

	// Dead letter queue for events not successfully sent to Lambda after the maximum retry period.
	dlq := awssqs.NewQueue(stack, jsii.String("cdk-stack-dead-letter-queue"), &awssqs.QueueProps{})

	// EventBridge rule to send events from the bus to Lambda function A.
	awsevents.NewRule(stack, jsii.String("cdk-stack-rule-lambda-function-a"), &awsevents.RuleProps{
		Description: jsii.String("Send all events in this bus to Lambda function A"),
		Enabled:     jsii.Bool(true),
		EventBus:    bus,
		RuleName:    jsii.String("cdk-stack-rule-lambda-function-a"),
		EventPattern: &awsevents.EventPattern{
			Region: &[]*string{jsii.String("eu-west-1")},
		},

		// Send events to Lambda or DQL.
		Targets: &[]awsevents.IRuleTarget{
			awseventstargets.NewLambdaFunction(fn_a, &awseventstargets.LambdaFunctionProps{
				DeadLetterQueue: dlq,
			})},
	})

	// EventBridge rule to send events from the bus to Lambda function B.
	awsevents.NewRule(stack, jsii.String("cdk-stack-rule-lambda-function-b"), &awsevents.RuleProps{
		Description: jsii.String("Send all events in this bus to Lambda function B"),
		Enabled:     jsii.Bool(true),
		EventBus:    bus,
		RuleName:    jsii.String("cdk-stack-rule-lambda-function-b"),
		EventPattern: &awsevents.EventPattern{
			Region: &[]*string{jsii.String("eu-west-1")},
		},

		// Send events to Lambda or DQL.
		Targets: &[]awsevents.IRuleTarget{
			awseventstargets.NewLambdaFunction(fn_b, &awseventstargets.LambdaFunctionProps{
				DeadLetterQueue: dlq,
			})},
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewEventsDrivenArchStack(app, "EventsDrivenArchStack", &EventsDrivenArchStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
