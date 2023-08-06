import { TodoResources } from "../lambda/todo-resources";
import { IRestApi } from "aws-cdk-lib/aws-apigateway";
import * as apigateway from 'aws-cdk-lib/aws-apigateway'

export class TodoApi {

  constructor(api: IRestApi, resources: TodoResources, authorizer: apigateway.IAuthorizer) {
    const todos = api.root.addResource("todos")

    // GET /todos
    todos.addMethod("GET", new apigateway.LambdaIntegration(resources.listHandler),{
      authorizer: authorizer
    })

    // GET /todos/:id
    const showTodo = todos.addResource("{id}")
    showTodo.addMethod("GET", new apigateway.LambdaIntegration(resources.getHandler), {
      authorizer: authorizer
    })

    // POST /todos/:id./_done
    const domeTodo = showTodo.addResource("_done").addMethod("POST", new apigateway.LambdaIntegration(resources.donetHandler), {
      authorizer: authorizer
    })

    // POST /todos
    todos.addMethod("POST", new apigateway.LambdaIntegration(resources.addHandler), {
      authorizer: authorizer
    })
  }

}