import { Construct } from "constructs";
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import * as iam from 'aws-cdk-lib/aws-iam';
import * as ddb from 'aws-cdk-lib/aws-dynamodb'

export class TodoResources {

  readonly todoTable: ddb.Table
  readonly todoCounterTable: ddb.Table
  readonly listHandler: lambda.Function
  readonly addHandler: lambda.Function
  readonly getHandler: lambda.Function
  readonly doneHandler: lambda.Function
  readonly deleteHandler: lambda.Function

  constructor(scope: Construct, todoTable: ddb.Table, todoCounterTable: ddb.Table) {
    this.todoTable = todoTable
    this.todoCounterTable = todoCounterTable
    this.listHandler = this.createListHandler(scope)
    this.addHandler = this.createAddHandler(scope)
    this.getHandler = this.createGetHandler(scope)
    this.doneHandler = this.createDoneHandler(scope)
    this.deleteHandler = this.createDeleteHandler(scope)
  }

  private createListHandler: (scope: Construct) => lambda.Function = scope => {
    const listHandler = new GoFunction(scope, 'todos-list-lambda', {
      entry: '../lambda/cmd/todos/list',
    })

    listHandler.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: [this.todoTable.tableArn],
      actions: [
          "dynamodb:Query",
      ]
    }))
    return listHandler
  }

  private createAddHandler: (scope: Construct) => lambda.Function = scope => {
    const addHandler = new GoFunction(scope, 'todos-add-lambda', {
      entry: '../lambda/cmd/todos/add',
    })
    addHandler.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: [this.todoTable.tableArn],
      actions: [
        "dynamodb:PutItem",
      ]
    }))
    addHandler.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: [this.todoCounterTable.tableArn],
      actions: [
        "dynamodb:UpdateItem",
      ]
    }))
    return addHandler
  }

  private createGetHandler: (scope: Construct) => lambda.Function = scope => {
    const addHandler = new GoFunction(scope, 'todos-get-lambda', {
      entry: '../lambda/cmd/todos/show',
    })
    addHandler.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: [this.todoTable.tableArn],
      actions: [
        "dynamodb:GetItem",
      ]
    }))
    return addHandler  
  }

  private createDoneHandler: (scope: Construct) => lambda.Function = scope => {
    const doneHandler = new GoFunction(scope, 'todos-done-lambda', {
      entry: '../lambda/cmd/todos/done',
    })
    doneHandler.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: [this.todoTable.tableArn],
      actions: [
        "dynamodb:UpdateItem",
      ]
    }))
    return doneHandler
  }

  private createDeleteHandler: (scope: Construct) => lambda.Function = scope => {
    const f = new GoFunction(scope, 'todos-delete-lambda', {
      entry: '../lambda/cmd/todos/delete',
    })
    f.addToRolePolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      resources: [this.todoTable.tableArn],
      actions: [
        "dynamodb:DeleteItem",
      ]
    }))
    return f
  }

}