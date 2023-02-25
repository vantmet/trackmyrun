package main

import (
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

	//To Use an existing image.
	imageName := jsii.String("606662134411.dkr.ecr.eu-west-2.amazonaws.com/trackmyrun:latest")

	//Create Fargate Service
	//Task Execution Role
	ter := awsiam.NewRole(stack, jsii.String("taskExecutionRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})
	ter.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ecr:BatchCheckLayerAvailability", "ecr:GetDownloadUrlForLayer", "ecr:BatchGetImage", "logs:CreateLogStream", "logs:PutLogEvents", "ecr:GetAuthorizationToken"),
		Resources: jsii.Strings("*"),
	}))
	//Task Role
	tr := awsiam.NewRole(stack, jsii.String("taskRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})
	tr.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents"),
		Resources: jsii.Strings("*"),
	}))
	td := awsecs.NewFargateTaskDefinition(stack, jsii.String("taskDefinition"), &awsecs.FargateTaskDefinitionProps{
		MemoryLimitMiB: jsii.Number(512),
		Cpu:            jsii.Number(256),
		ExecutionRole:  ter,
		TaskRole:       tr,
	})

	container := td.AddContainer(jsii.String("taskContainer"), &awsecs.ContainerDefinitionOptions{
		Image: awsecs.ContainerImage_FromRegistry(imageName, &awsecs.RepositoryImageProps{}),
		Logging: awsecs.LogDriver_AwsLogs(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("task"),
		}),
	})

	container.AddPortMappings(&awsecs.PortMapping{ContainerPort: jsii.Number(5000)})

	service := awsecs.NewFargateService(stack, jsii.String("TMRFGService"), &awsecs.FargateServiceProps{
		Cluster:        cluster,
		TaskDefinition: td,
		DesiredCount:   jsii.Number(1),
		AssignPublicIp: jsii.Bool(true),
	})

	// Add and Application load-balancer
	targetGroup := awselasticloadbalancingv2.NewApplicationTargetGroup(stack, jsii.String("TRMALB"), &awselasticloadbalancingv2.ApplicationTargetGroupProps{
		Port:       jsii.Number(5000),
		Vpc:        vpc,
		Protocol:   awselasticloadbalancingv2.ApplicationProtocol_HTTP,
		TargetType: awselasticloadbalancingv2.TargetType_IP,
		HealthCheck: &awselasticloadbalancingv2.HealthCheck{
			Path: jsii.String("/ping"),
		},
	})
	lb := awselasticloadbalancingv2.NewApplicationLoadBalancer(stack, jsii.String("TMRLB"), &awselasticloadbalancingv2.ApplicationLoadBalancerProps{
		Vpc:            vpc,
		InternetFacing: jsii.Bool(true),
	})

	listener := lb.AddListener(jsii.String("Listener"), &awselasticloadbalancingv2.BaseApplicationListenerProps{
		Port: jsii.Number(80),
	})
	listener.AddTargetGroups(jsii.String("FargateVMs"), &awselasticloadbalancingv2.AddApplicationTargetGroupsProps{
		TargetGroups: &[]awselasticloadbalancingv2.IApplicationTargetGroup{targetGroup},
	})

	service.AttachToApplicationTargetGroup(targetGroup)

	awscdk.NewCfnOutput(stack, jsii.String("URL"), &awscdk.CfnOutputProps{
		Value:       lb.LoadBalancerDnsName(),
		Description: jsii.String("The URL of the load balancer for testing."),
		ExportName:  jsii.String("TMRURL"),
	})

	awscognito.NewUserPool(stack, jsii.String("tmruserpool"), &awscognito.UserPoolProps{
		UserPoolName: jsii.String("Track My Run - userpool"),
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
