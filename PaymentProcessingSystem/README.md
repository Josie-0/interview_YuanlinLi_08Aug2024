## Payment Processing System

### API file
[查看API说明](api.yaml)

### Postman collection
[查看request/response示例](https://api.postman.com/collections/37534429-6f843dc7-01b5-40bb-9478-f0d96c199363?access_key=PMAT-01J4W48TWE2F034QFBNYWCHSAJ)

### RUN

`docker-compose build`

`docker-compose up -d`

### 可优化方向
1. 输入验证： 对用户输入的数据进行严格验证
2. 考虑异步处理耗时长的支付 
3. 使用 OAuth 2.0 或 JWT（JSON Web Token）进行更安全的认证和授权。