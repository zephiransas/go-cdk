# Welcome to your CDK TypeScript project

This is a blank project for CDK development with TypeScript.

The `cdk.json` file tells the CDK Toolkit how to execute your app.

## Useful commands

* `npm run build`   compile typescript to js
* `npm run watch`   watch for changes and compile
* `npm run test`    perform the jest unit tests
* `npm run lint`    lint
* `cdk deploy`      deploy this stack to your default AWS account/region
* `cdk diff`        compare deployed stack with current state
* `cdk synth`       emits the synthesized CloudFormation template

## 事前に作っておくべきAWSリソース

### Cognito User Pool

TBD

### SSM Parameter Store

Cognito作成後に、SSM Parameter Storeに以下のパラメータを設定する

- `/go-cdk/poolId` - Cognito User PoolのPool ID（ARNではない）
- `/go-cdk/clientId` - アプリから利用するクライアントID
- `/go-cdk/clientSecret` - Secure Stringにすること。アプリから利用するクライアントのSecret

```
aws ssm put-parameter --name "/go-cdk/poolId" --value "REPLACE_HERE" --type String
aws ssm put-parameter --name "/go-cdk/clientId" --value "REPLACE_HERE" --type String
aws ssm put-parameter --name "/go-cdk/clientSecret" --value "REPLACE_HERE" --type SecureString
```