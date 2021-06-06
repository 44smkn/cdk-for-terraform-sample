# CDK FOR TERRAFORM SAMPLE

## Reference

* [Announcing CDK for Terraform 0.4](https://www.hashicorp.com/blog/announcing-cdk-for-terraform-0-4)
* [Getting Started with Go](https://github.com/hashicorp/terraform-cdk/blob/main/docs/getting-started/go.md)
* [Terraform CDK Example Go × AWS](https://github.com/hashicorp/terraform-cdk/tree/main/examples/go/aws)

## Runbook

```sh
cdktf init --template="go" --local
vi cdktf.json # Declare the dependent Terraform Providers and Modules
cdktf get
cdktf diff
cdktf destroy
```
