import * as cdk from 'aws-cdk-lib';
import * as ddb from 'aws-cdk-lib/aws-dynamodb'
import * as apigateway from 'aws-cdk-lib/aws-apigateway'
import * as lambda from '@aws-cdk/aws-lambda-go-alpha'
import {RemovalPolicy, aws_iam as iam} from 'aws-cdk-lib';
import {Construct} from 'constructs';

export class CdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Dynamodb Table
    const table = new ddb.Table(this, 'todo-table', {
      tableName: "todos-go",
      partitionKey: {
        name: "user_id",
        type: ddb.AttributeType.STRING
      },
      sortKey: {
        name: "id",
        type: ddb.AttributeType.NUMBER
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

    const addHandler = new lambda.GoFunction(this, 'add-lambda', {
      entry: '../lambda/cmd/add',
    })
    addHandler.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: [table.tableArn],
      actions: [
        "dynamodb:PutItem",
      ]
    }))

    // API Gateway
    const api = new apigateway.RestApi(this, "todo-api")

    const todos = api.root.addResource("todos")

    // GET /todos
    todos.addMethod("GET", new cdk.aws_apigateway.LambdaIntegration(listHandler))
    
    // POST /todos
    todos.addMethod("POST", new cdk.aws_apigateway.LambdaIntegration(addHandler))
  }
}
