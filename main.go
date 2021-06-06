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

	vpc := aws.NewVpc(stack, jsii.String("isucon_vpc"), &aws.VpcConfig{
		CidrBlock: jsii.String("172.16.0.0/16"),
		Tags: &map[string]*string{
			"Name": jsii.String("isucon-training"),
		},
	})

	igw := aws.NewInternetGateway(stack, jsii.String("isucon_vpc_igw"), &aws.InternetGatewayConfig{
		VpcId: vpc.Id(),
	})

	subnet := aws.NewSubnet(stack, jsii.String("isucon9_subnet"), &aws.SubnetConfig{
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

	// routeTableでCidr指定だと。 "destination_prefix_list_id" is required というエラーが発生したため
	managedPrefix := aws.NewEc2ManagedPrefixList(stack, jsii.String("default_prefix"), &aws.Ec2ManagedPrefixListConfig{
		Name:          jsii.String("default CIDR"),
		AddressFamily: jsii.String("IPv4"),
		MaxEntries:    jsii.Number(1),
		Entry: &[]*aws.Ec2ManagedPrefixListEntry{
			{
				Cidr:        jsii.String("0.0.0.0/0"),
				Description: jsii.String("default"),
			},
		},
	})

	routeTable := aws.NewDefaultRouteTable(stack, jsii.String("isucon9_subnet_route_table"), &aws.DefaultRouteTableConfig{
		DefaultRouteTableId: vpc.DefaultRouteTableId(),
		Route: &[]*aws.DefaultRouteTableRoute{
			{
				GatewayId:               igw.Id(),
				DestinationPrefixListId: managedPrefix.Id(),
			},
		},
		DependsOn: &[]cdktf.ITerraformDependable{
			igw,
		},
	})

	aws.NewRouteTableAssociation(stack, jsii.String("isucon9_subnet_route_table_d"), &aws.RouteTableAssociationConfig{
		SubnetId:     subnet.Id(),
		RouteTableId: routeTable.Id(),
		DependsOn: &[]cdktf.ITerraformDependable{
			igw,
		},
	})

	keyPair := aws.NewKeyPair(stack, jsii.String("developer_keypair"), &aws.KeyPairConfig{
		KeyName:   jsii.String("developer-key"),
		PublicKey: jsii.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCigWIwCDN7ME+WVisrQzZVEvjfRPkhkMSdJjnz0a4SvaCvbCVPXJ/xZpqAtmnT8gir0s+MaAdyAYEmqbnoI/m30vuINRQ9w8tFUSjuw/9fQO2wZHpoS+PvPpM1XnOeGpLLf3zmjVOvyPOooO9ufs+W+0bZPb2emhaViUqflMsqQ8a2O2gc7jH1bktAeb7/KwUCMJAs/m6/07kNgQ/fJ3LaJF2w7vg1KQq831gZxR6mbBJIttuctymqJcUGwxBRWTxi19GgQdz4og+DjaKBsdjaBn1UYs0YB304iImULXouWo6caK1+i7eohS1vPIWVJ6iCSN8Qs51ofGIC22+1iSwo/GBiePHiBnqM9JcC0wBDCKqJdxTnLc12eU2T2djd/weYmkQOzJ7AeDufHToEF/fV8GzZcveQosBtTamat2UmhzhpKHUdtMOaDpXq3cchmdx0f7wE7u/YuAgEny64LgF4zmqo25OEh+cYRcbUFi0w5ny1j5jqwiPHtlpFLKLAO4jQWfNFRPVLcgEUhMSzkNh3O6dfRnKByUsj4AU+3ajuyehoncrtKZyziPEQfmbAbMmHBar0KiXvVw1m++3F2AcsaYBZ4HvRlQ2uS2HjvtiMs5si+sU8jC8jF+dLp6wDNsw2z7vaJ1sm1bUmLJa3gQFVA29c0/EazKdcVgFxfakAHw== jrywm121@gmail.com"),
	})

	securityGroup := aws.NewSecurityGroup(stack, jsii.String("isucon9_qualify_instance_sg"), &aws.SecurityGroupConfig{
		Name:  jsii.String("allow_ssh"),
		VpcId: vpc.Id(),
		Ingress: &[]*aws.SecurityGroupIngress{
			{
				FromPort:       jsii.Number(22),
				ToPort:         jsii.Number(22),
				Protocol:       jsii.String("tcp"),
				CidrBlocks:     jsii.Strings("0.0.0.0/0"),
				Ipv6CidrBlocks: jsii.Strings("::/0"),
			},
			{
				FromPort:       jsii.Number(8000),
				ToPort:         jsii.Number(8000),
				Protocol:       jsii.String("tcp"),
				CidrBlocks:     jsii.Strings("0.0.0.0/0"),
				Ipv6CidrBlocks: jsii.Strings("::/0"),
			},
		},
		Egress: &[]*aws.SecurityGroupEgress{
			{
				FromPort:       jsii.Number(0),
				ToPort:         jsii.Number(0),
				Protocol:       jsii.String("-1"),
				CidrBlocks:     jsii.Strings("0.0.0.0/0"),
				Ipv6CidrBlocks: jsii.Strings("::/0"),
			},
		},
	})

	privateIp := jsii.String("172.16.10.100")
	instace := aws.NewInstance(stack, jsii.String("isucon9_qualify_instance"), &aws.InstanceConfig{
		Ami:                 jsii.String("ami-03b1b78bb1da5122f"),
		InstanceType:        jsii.String("c5.large"),
		KeyName:             keyPair.KeyName(),
		VpcSecurityGroupIds: &[]*string{securityGroup.Id()},

		PrivateIp: privateIp,
		SubnetId:  subnet.Id(),
	})

	aws.NewEip(stack, jsii.String("isucon9_qualify_instance_eip"), &aws.EipConfig{
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
