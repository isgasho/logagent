# logagent
前端监控-异常,api调用

**后端:**
- 安装启动nginx,将logagent/web/log.gif放到www访问目录,地址可自定

**优点：**
- 解决跨域问题
- 不会影响过多占用业务资源
- get请求的url自带缓存，避免重复请求，过度浪费和造成太大压力
- kafka消峰
- offset读取

**缺点**
- url拼接参数有长度大小限制
- url只支持url编码,很容易被解析出来,不能传输敏感数据

**解析格式** <br/>
http://www.hank.com/log.gif?log={"aa":"aa","bb":"bb"} <br/>
解析data参数后直接json解析成对象

#未完善
配置文件管理比较混乱
elastic请求过程优化，没发送一次记录都要初始化资源上传

未解决:
文件读取的特殊字符问题



