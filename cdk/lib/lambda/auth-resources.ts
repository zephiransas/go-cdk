import { Construct } from "constructs";
import * as lambda from 'aws-cdk-lib/aws-lambda'
import { StringParameter } from "aws-cdk-lib/aws-ssm";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import { Effect, PolicyStatement } from "aws-cdk-lib/aws-iam";

export class AuthResources {

  readonly loginHandler: lambda.Function
  readonly authorizeHandler: lambda.Function

  constructor(scope: Construct) {
    this.loginHandler = this.createLoginHandler(scope)
    this.authorizeHandler = this.createAuthorizeHandler(scope)
  }

  private createLoginHandler: (scope: Construct) => lambda.Function = scope => {

    const poolId = StringParameter.valueForStringParameter(scope, "/go-cdk/poolId")
    const clientId = StringParameter.valueForStringParameter(scope, "/go-cdk/clientId")

    const loginHandler = new GoFunction(scope, 'login-lambda', {
      entry: '../lambda/cmd/login',
      environment: {
        "POOL_ID": poolId,
        "CLIENT_ID": clientId,
        // CLIENT_SECRETは、都度Lambdaから取得する
      }
    })
    loginHandler.addToRolePolicy(new PolicyStatement({
      effect: Effect.ALLOW,
      resources: [`arn:aws:cognito-idp:ap-northeast-1:${process.env.AWS_ACCOUNT_ID}:userpool/${process.env.USER_POOL_ID}`],
      actions: [
        "cognito-idp:AdminInitiateAuth",
      ]
    }))
    loginHandler.addToRolePolicy(new PolicyStatement({
      effect: Effect.ALLOW,
      resources: [`arn:aws:ssm:ap-northeast-1:${process.env.AWS_ACCOUNT_ID}:parameter/go-cdk/clientSecret`],
      actions: [
        "ssm:GetParameter",
      ]
    }))

    return loginHandler
  }

  private createAuthorizeHandler: (scope: Construct) => lambda.Function = scope => {
    return new GoFunction(scope, 'authorizer-lambda', {
      entry: '../lambda/cmd/authorizer'
    })
  }

}