# CDK FOR TERRAFORM SAMPLE

## Reference

* [Announcing CDK for Terraform 0.4](https://www.hashicorp.com/blog/announcing-cdk-for-terraform-0-4)
* [Getting Started with Go](https://github.com/hashicorp/terraform-cdk/blob/main/docs/getting-started/go.md)
* [Terraform CDK Example Go × AWS](https://github.com/hashicorp/terraform-cdk/tree/main/examples/go/aws)

## Runbook

```sh
cdktf init --template="go" --local
vi cdktf.json # 依存するTerraformのProviderとModuleを宣言
cdktf get # codeがgenerateされる
go mod tidy
```

