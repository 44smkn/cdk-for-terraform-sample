package main

import (
	"github.com/44smkn/cdk-for-terraform-sample/generated/hashicorp/aws"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

const (
	region = "ap-northeast-1"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	aws.NewAwsProvider(stack, jsii.String("aws"), &aws.AwsProviderConfig{
		Region: jsii.String(region),
	})

	vpc := aws.NewVpc(scope, jsii.String("isucon_vpc"), &aws.VpcConfig{
		CidrBlock: jsii.String("172.16.0.0/16"),
		Tags: &map[string]*string{
			"Name": jsii.String("isucon-training"),
		},
	})

	igw := aws.NewInternetGateway(scope, jsii.String("isucon_vpc_igw"), &aws.InternetGatewayConfig{
		VpcId: vpc.Id(),
	})

	subnet := aws.NewSubnet(scope, jsii.String("isucon9_subnet"), &aws.SubnetConfig{
		VpcId:            vpc.Id(),
		CidrBlock:        jsii.String("172.16.10.0/24"),
		AvailabilityZone: jsii.String("ap-northeast-1d"),
		Tags: &map[string]*string{
			"Name": jsii.String("isucon-training"),
		},
		DependsOn: &[]cdktf.ITerraformDependable{
			igw,
		},
	})

	keyPair := aws.NewKeyPair(scope, jsii.String("developer"), &aws.KeyPairConfig{
		KeyName:   jsii.String("developer-key"),
		PublicKey: jsii.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCigWIwCDN7ME+WVisrQzZVEvjfRPkhkMSdJjnz0a4SvaCvbCVPXJ/xZpqAtmnT8gir0s+MaAdyAYEmqbnoI/m30vuINRQ9w8tFUSjuw/9fQO2wZHpoS+PvPpM1XnOeGpLLf3zmjVOvyPOooO9ufs+W+0bZPb2emhaViUqflMsqQ8a2O2gc7jH1bktAeb7/KwUCMJAs/m6/07kNgQ/fJ3LaJF2w7vg1KQq831gZxR6mbBJIttuctymqJcUGwxBRWTxi19GgQdz4og+DjaKBsdjaBn1UYs0YB304iImULXouWo6caK1+i7eohS1vPIWVJ6iCSN8Qs51ofGIC22+1iSwo/GBiePHiBnqM9JcC0wBDCKqJdxTnLc12eU2T2djd/weYmkQOzJ7AeDufHToEF/fV8GzZcveQosBtTamat2UmhzhpKHUdtMOaDpXq3cchmdx0f7wE7u/YuAgEny64LgF4zmqo25OEh+cYRcbUFi0w5ny1j5jqwiPHtlpFLKLAO4jQWfNFRPVLcgEUhMSzkNh3O6dfRnKByUsj4AU+3ajuyehoncrtKZyziPEQfmbAbMmHBar0KiXvVw1m++3F2AcsaYBZ4HvRlQ2uS2HjvtiMs5si+sU8jC8jF+dLp6wDNsw2z7vaJ1sm1bUmLJa3gQFVA29c0/EazKdcVgFxfakAHw== jrywm121@gmail.com"),
	})

	privateIp := jsii.String("172.16.10.100")
	instace := aws.NewInstance(scope, jsii.String("isucon9_qualify_instance"), &aws.InstanceConfig{
		Ami:          jsii.String("ami-03b1b78bb1da5122f"),
		InstanceType: jsii.String("c5.large"),
		KeyName:      keyPair.KeyName(),

		PrivateIp: privateIp,
		SubnetId:  subnet.Id(),
	})

	aws.NewEip(scope, jsii.String("isucon9_qualify_instance"), &aws.EipConfig{
		Vpc:                    jsii.Bool(true),
		Instance:               instace.Id(),
		AssociateWithPrivateIp: privateIp,
		DependsOn: &[]cdktf.ITerraformDependable{
			igw,
		},
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	NewMyStack(app, "cdk-for-terraform-sample")

	app.Synth()
}
