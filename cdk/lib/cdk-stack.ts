import * as cdk from 'aws-cdk-lib';
import {RemovalPolicy, aws_iam as iam} from 'aws-cdk-lib';
import {Construct} from 'constructs';
import * as lambda from '@aws-cdk/aws-lambda-go-alpha'
import {AttributeType} from "aws-cdk-lib/aws-dynamodb";

export class CdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const table = new cdk.aws_dynamodb.Table(this, 'todo-table', {
      tableName: "todos-go",
      partitionKey: {
        name: "user_id",
        type: AttributeType.STRING
      },
      sortKey: {
        name: "id",
        type: AttributeType.NUMBER
      },
      readCapacity: 1,
      writeCapacity: 1,
      removalPolicy: RemovalPolicy.RETAIN
    })

    const listHandler = new lambda.GoFunction(this, 'list-lambda', {
      entry: '../lambda/cmd/list',
    })
    listHandler.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: [table.tableArn],
      actions: [
          "dynamodb:Query",
      ]
    }))

    const api = new cdk.aws_apigateway.RestApi(this, "todo-api")

    const books = api.root.addResource("todos")
    books.addMethod("GET", new cdk.aws_apigateway.LambdaIntegration(listHandler))

  }
}
