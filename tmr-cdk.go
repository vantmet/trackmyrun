package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type TmrCdkStackProps struct {
	awscdk.StackProps
}

func NewTmrCdkStack(scope constructs.Construct, id string, props *TmrCdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	userPool := awscognito.NewUserPool(stack, jsii.String("tmruserpool"), &awscognito.UserPoolProps{
		UserPoolName: jsii.String("Track My Run - userpool"),
	})
	userPool.ApplyRemovalPolicy(awscdk.RemovalPolicy_DESTROY)
	userPoolClient := userPool.AddClient(jsii.String("TMR Users"), &awscognito.UserPoolClientOptions{})
	clientID := userPoolClient.UserPoolClientId()
	poolID := userPool.UserPoolId()

	//Create a VPC + Cluster to live in
	vpc := awsec2.NewVpc(stack, jsii.String("TMRVPC"), &awsec2.VpcProps{})
	cluster := awsecs.NewCluster(stack, jsii.String("TMRCluster"), &awsecs.ClusterProps{
		Vpc: vpc,
	})

	// To create the Docker image.
	/*
		appImage := awsecrassets.NewDockerImageAsset(stack, jsii.String("ApplicationImage"), &awsecrassets.DockerImageAssetProps{
			Directory: jsii.String("./app"),
		})
		imageName := jsii.String(*appImage.ImageUri()) */

	// Add target groups
	appTargetGroup := awselasticloadbalancingv2.NewApplicationTargetGroup(stack, jsii.String("TRMALB"), &awselasticloadbalancingv2.ApplicationTargetGroupProps{
		Port:       jsii.Number(5000),
		Vpc:        vpc,
		Protocol:   awselasticloadbalancingv2.ApplicationProtocol_HTTP,
		TargetType: awselasticloadbalancingv2.TargetType_IP,
		HealthCheck: &awselasticloadbalancingv2.HealthCheck{
			Path: jsii.String("/ping"),
		},
	})
	authTargetGroup := awselasticloadbalancingv2.NewApplicationTargetGroup(stack, jsii.String("TRMAuthALB"), &awselasticloadbalancingv2.ApplicationTargetGroupProps{
		Port:       jsii.Number(5000),
		Vpc:        vpc,
		Protocol:   awselasticloadbalancingv2.ApplicationProtocol_HTTP,
		TargetType: awselasticloadbalancingv2.TargetType_IP,
		HealthCheck: &awselasticloadbalancingv2.HealthCheck{
			Path: jsii.String("/ping"),
		},
	})

	//To Use an existing image.
	appImageName := jsii.String("606662134411.dkr.ecr.eu-west-2.amazonaws.com/trackmyrun:latest")
	authImageName := jsii.String("606662134411.dkr.ecr.eu-west-2.amazonaws.com/trackmyrun-auth:latest")

	//Create Fargate Service
	//Task Execution Role
	taskExecutionRole := awsiam.NewRole(stack, jsii.String("taskExecutionRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})
	taskExecutionRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ecr:BatchCheckLayerAvailability", "ecr:GetDownloadUrlForLayer", "ecr:BatchGetImage", "logs:CreateLogStream", "logs:PutLogEvents", "ecr:GetAuthorizationToken"),
		Resources: jsii.Strings("*"),
	}))
	//Task Role
	taskRole := awsiam.NewRole(stack, jsii.String("taskRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})
	taskRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents"),
		Resources: jsii.Strings("*"),
	}))
	appTaskDefinition := awsecs.NewFargateTaskDefinition(stack, jsii.String("appTaskDefinition"), &awsecs.FargateTaskDefinitionProps{
		MemoryLimitMiB: jsii.Number(512),
		Cpu:            jsii.Number(256),
		ExecutionRole:  taskExecutionRole,
		TaskRole:       taskRole,
	})

	appContainer := appTaskDefinition.AddContainer(jsii.String("taskContainer"), &awsecs.ContainerDefinitionOptions{
		Image: awsecs.ContainerImage_FromRegistry(appImageName, &awsecs.RepositoryImageProps{}),
		Logging: awsecs.LogDriver_AwsLogs(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("task"),
		}),
	})

	appContainer.AddPortMappings(&awsecs.PortMapping{ContainerPort: jsii.Number(5000)})

	authTaskDefinition := awsecs.NewFargateTaskDefinition(stack, jsii.String("authTaskDefinition"), &awsecs.FargateTaskDefinitionProps{
		MemoryLimitMiB: jsii.Number(512),
		Cpu:            jsii.Number(256),
		ExecutionRole:  taskExecutionRole,
		TaskRole:       taskRole,
	})

	authContainer := authTaskDefinition.AddContainer(jsii.String("taskContainer"), &awsecs.ContainerDefinitionOptions{
		Image: awsecs.ContainerImage_FromRegistry(authImageName, &awsecs.RepositoryImageProps{}),
		Logging: awsecs.LogDriver_AwsLogs(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("task"),
		}),
	})
	authContainerPort := float64(5001)
	authContainerPortString := fmt.Sprintf("%f", authContainerPort)

	authContainer.AddPortMappings(&awsecs.PortMapping{ContainerPort: jsii.Number(authContainerPort)})
	authContainer.AddEnvironment(jsii.String("PORT"), jsii.String(authContainerPortString))
	authContainer.AddEnvironment(jsii.String("COGNITO_APP_CLIENT_ID"), jsii.String(*clientID))
	authContainer.AddEnvironment(jsii.String("COGNITO_APP_POOL_ID"), jsii.String(*poolID))

	appService := awsecs.NewFargateService(stack, jsii.String("TMR Main Service"), &awsecs.FargateServiceProps{
		Cluster:        cluster,
		TaskDefinition: appTaskDefinition,
		DesiredCount:   jsii.Number(1),
		AssignPublicIp: jsii.Bool(true),
	})

	appService.AttachToApplicationTargetGroup(appTargetGroup)

	authService := awsecs.NewFargateService(stack, jsii.String("TMR Auth Service"), &awsecs.FargateServiceProps{
		Cluster:        cluster,
		TaskDefinition: authTaskDefinition,
		DesiredCount:   jsii.Number(1),
		AssignPublicIp: jsii.Bool(true),
	})

	authService.AttachToApplicationTargetGroup(authTargetGroup)

	lb := awselasticloadbalancingv2.NewApplicationLoadBalancer(stack, jsii.String("TMRLB"), &awselasticloadbalancingv2.ApplicationLoadBalancerProps{
		Vpc:            vpc,
		InternetFacing: jsii.Bool(true),
	})

	applistener := lb.AddListener(jsii.String("appListener"), &awselasticloadbalancingv2.BaseApplicationListenerProps{
		Port: jsii.Number(80),
	})
	applistener.AddTargetGroups(jsii.String("appVMs"), &awselasticloadbalancingv2.AddApplicationTargetGroupsProps{
		TargetGroups: &[]awselasticloadbalancingv2.IApplicationTargetGroup{appTargetGroup},
	})

	authlistener := lb.AddListener(jsii.String("authListener"), &awselasticloadbalancingv2.BaseApplicationListenerProps{
		Port:     jsii.Number(8080),
		Protocol: awselasticloadbalancingv2.ApplicationProtocol_HTTP,
	})

	authlistener.AddTargetGroups(jsii.String("authVMs"), &awselasticloadbalancingv2.AddApplicationTargetGroupsProps{
		TargetGroups: &[]awselasticloadbalancingv2.IApplicationTargetGroup{authTargetGroup},
	})

	awscdk.NewCfnOutput(stack, jsii.String("URL"), &awscdk.CfnOutputProps{
		Value:       lb.LoadBalancerDnsName(),
		Description: jsii.String("The URL of the load balancer for testing."),
		ExportName:  jsii.String("TMRURL"),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewTmrCdkStack(app, "TmrCdkStack", &TmrCdkStackProps{
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
