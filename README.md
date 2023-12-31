# What's this?

- サーバレスでTODO APIを構成するためのサンプル

## 環境

- AWS
    - API Gateway
    - Lambda
    - DynamoDB
    - Cognito User Pool
    - CDK
- Golang

## 事前に作っておくべきAWSリソース

### Cognito User Pool

以下のコマンドで適当なユーザを作成

```
aws cognito-idp admin-create-user --user-pool-id [ユーザプールID] \
--username [一意なユーザ名] \
--user-attributes Name=email,Value="hogefuga@example.com" Name=email_verified,Value=true \
--message-action SUPPRESS
```

その後、以下のコマンドでパスワードを設定

```
aws cognito-idp admin-set-user-password \
--user-pool-id [ユーザプールID] \
--username [一意なユーザ名] \
--password [パスワード] \
--permanent
```

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