import * as cdk from 'aws-cdk-lib';
import * as ddb from 'aws-cdk-lib/aws-dynamodb'
import * as apigateway from 'aws-cdk-lib/aws-apigateway'
import * as lambda from '@aws-cdk/aws-lambda-go-alpha'
import * as logs from 'aws-cdk-lib/aws-logs'
import {RemovalPolicy, aws_iam as iam} from 'aws-cdk-lib';
import {Construct} from 'constructs';
import { StringParameter } from 'aws-cdk-lib/aws-ssm';

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


    const poolId = StringParameter.valueForStringParameter(this, "/go-cdk/poolId")
    const clientId = StringParameter.valueForStringParameter(this, "/go-cdk/clientId")

    const loginHandler = new lambda.GoFunction(this, 'login-lambda', {
      entry: '../lambda/cmd/login',
      environment: {
        "POOL_ID": poolId,
        "CLIENT_ID": clientId,
        // CLIENT_SECRETは、都度Lambdaから取得する
      }
    })
    loginHandler.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: ["arn:aws:cognito-idp:ap-northeast-1:919951165082:userpool/ap-northeast-1_43WZ6LiP3"],
      actions: [
        "cognito-idp:AdminInitiateAuth",
      ]
    }))
    loginHandler.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: ["arn:aws:ssm:ap-northeast-1:919951165082:parameter/go-cdk/clientSecret"],
      actions: [
        "ssm:GetParameter",
      ]
    }))


    const authorizerHandler = new lambda.GoFunction(this, 'authorizer-lambda', {
      entry: '../lambda/cmd/authorizer'
    })

    const authorizer = new apigateway.TokenAuthorizer(this, 'token-authorizer', {
      handler: authorizerHandler,
      identitySource: apigateway.IdentitySource.header("Authorization")
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
    addHandler.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: [counterTable.tableArn],
      actions: [
        "dynamodb:UpdateItem",
      ]
    }))

    // API Gateway
    const api = new apigateway.RestApi(this, "todo-api")

    const oauth = api.root.addResource("oauth")
    oauth.addResource("login").addMethod("GET", new apigateway.LambdaIntegration(loginHandler))

    const todos = api.root.addResource("todos")

    // GET /todos
    todos.addMethod("GET", new cdk.aws_apigateway.LambdaIntegration(listHandler),{
      authorizer: authorizer
    })
    
    // POST /todos
    todos.addMethod("POST", new cdk.aws_apigateway.LambdaIntegration(addHandler), {
      authorizer: authorizer
    })
  }
}
