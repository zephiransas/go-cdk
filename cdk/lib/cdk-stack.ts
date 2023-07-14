import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from '@aws-cdk/aws-lambda-go-alpha'

export class CdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new lambda.GoFunction(this, 'handler', {
      entry: '../lambda/cmd/list'
    })

  }
}
