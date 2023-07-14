import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from '@aws-cdk/aws-lambda-go-alpha'

export class CdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const handler = new lambda.GoFunction(this, 'handler', {
      entry: '../lambda/cmd/list'
    })

    const api = new cdk.aws_apigateway.RestApi(this, "todo-api")

    const books = api.root.addResource("todos")
    books.addMethod("GET", new cdk.aws_apigateway.LambdaIntegration(handler))

  }
}
