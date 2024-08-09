## Game Room Management System

### API file
[查看API说明](api.yaml)

### Postman collection
[查看request/response示例](https://api.postman.com/collections/37534429-2787df32-be03-43c0-9e35-8e3421477c0c?access_key=PMAT-01J4W46RV8CN0JZRTC4MNTPR5R)

### RUN

`docker-compose build`

`docker-compose up -d`

### 可优化方向
1. 支持批量操作
2. 输入验证：将请求的验证逻辑提取到单独的验证函数或使用专门的验证库
3. 考虑使用中间件来统一处理错误响应。这样可以避免在每个处理函数中重复相同的错误处理逻辑。
4. 数据库查询优化：
   1. 对于查询返回大量数据的接口（如预订记录），考虑添加分页功能