# web-monitor
前端监控-异常,api调用

**后端:**
- 安装启动nginx,将web-monitor/web/log.gif放到www访问目录,地址可自定


**优点：**
- 最重要还是能解决跨域问题
- 不会影响过多占用业务资源
- get请求的url自带缓存，避免重复请求，过度浪费和造成太大压力
- 
**缺点**
- url拼接参数有长度大小限制
- url只支持url编码,很容易被解析出来,不能传输敏感数据

**解析格式两种** <br/>

(1)http://www.hank.com/log.gif?aa="aa"&bb="bb" <br/>
可以通过url参数解析出参数出来

(2)http://www.hank.com/log.gif?log={"aa":"aa","bb":"bb"} <br/>
解析data参数后直接json解析成对象



