import * as cdk from 'aws-cdk-lib';
import * as ddb from 'aws-cdk-lib/aws-dynamodb'
import * as apigateway from 'aws-cdk-lib/aws-apigateway'
import {RemovalPolicy, aws_iam as iam} from 'aws-cdk-lib';
import {Construct} from 'constructs';
import { TodoResources } from './lambda/todo-resources';
import { AuthResources } from './lambda/auth-resources';
import { TodoApi } from './api/todo-api';

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
      removalPolicy: RemovalPolicy.DESTROY  //WARN
    })

    const counterTable = new ddb.Table(this, 'todo-counter-table', {
      tableName: "todos-go-counter",
      partitionKey: {
        name: "user_id",
        type: ddb.AttributeType.STRING
      },
      readCapacity: 1,
      writeCapacity: 1,
      removalPolicy: RemovalPolicy.DESTROY  //WARN
    })

    const authResources = new AuthResources(this)

    const authorizer = new apigateway.TokenAuthorizer(this, 'token-authorizer', {
      handler: authResources.authorizeHandler,
      resultsCacheTtl: cdk.Duration.seconds(0),   //cacheを無効にする
      identitySource: apigateway.IdentitySource.header("Authorization")
    })

    const todoResources = new TodoResources(this, table, counterTable)

    // API Gateway
    const api = new apigateway.RestApi(this, "todo-api")

    const oauth = api.root.addResource("oauth")
    oauth.addResource("login").addMethod("GET", new apigateway.LambdaIntegration(authResources.loginHandler))


    // Todo
    new TodoApi(api, todoResources, authorizer)

  }
}
