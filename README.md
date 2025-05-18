## 内容总结：

1. 数据库的三个时间戳是自动维护的，可以保留
2. gorm默认软删除且支持查询标记删除的数据，不用手动维护状态
3. vnet表的补充ID，和name完全相同
4. Nickname没有删，生成的关于用户profile的内容用到nickname的太多了，先放在那里
5. 返回总页数直接在GetUsagePageResponseData里加了一个PageCount
6. 删除所有生成的nickname相关
7. 修改并完善用户Profile相关api，补全IsAdmin、IsVip
