## Endless Challenge System

### API file
[查看API说明](api.yaml)

### Postman collection
[查看request/response示例](https://api.postman.com/collections/37534429-29dc2b97-fc7b-4909-9780-842acdd1c2b7?access_key=PMAT-01J4W3PHHX0XCEDNMT6B863D9R)

### RUN

`docker-compose build`

`docker-compose up -d`

### 可优化方向
1. 当前使用内存中的数据结构来存储挑战信息，可能会出现并发访问问题。可以考虑加入数据库/缓存，频繁查询的数据可放在缓存
2. 考虑引入挑战处理的异步任务队列
3. 对挑战请求的频率进行限制