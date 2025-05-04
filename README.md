## 待完成:

1. 还没去看数据库连接和迁移部分，usage和vnet的repository未连接

## 内容总结：

1. Model.User：添加Username项
2. 修改RegisterRequest，添加Username和NickName项，均为非空
3. 将LoginRequest的email改成EmailOrUsername来实现邮箱或用户名登录
4. Service.User：Login添加用户名登录方式（优先级低于邮箱登录），未实现邮箱检测，按邮箱-用户名顺序查询数据库
5. 新增Usage相关内容：Model.Usage，Repository.Usage，实现service层的GetUsage（按近30天，近7天，近12月）其中GetUsageRequest的range设定为“month”、“30days”或“7days”，如果为month返回数据列表为每月总量，否则为日总量
6. 实现vnet相关内容，其中model.vnet内容为get vnet的Response返回内容加上UserId
7. 没有对DELETE vnet时删除内容不存在做特殊处理
8. admin和user功能相似的部分直接复制过去修改的代码，考虑到拓展性没有直接调用
9. 最后的/user/vnet/<USERID>/<VNETID>是不是写错了，是admin吧
10. admin的vnet的删改部分没有用到USERID，直接用的VNETID找

## 奇怪的问题：

1. 邮箱和用户名：
   例子：A用户的username叫做“a@b.com”，而a@b.com刚好是B用户的注册邮箱怎么处理
