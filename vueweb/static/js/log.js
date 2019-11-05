const LogLevel = {
  LevelEmergency: 0,
  LevelAlert: 1,
  LevelCritical: 2,
  LevelError: 3,
  LevelWarning: 4,
  LevelNotice: 5,
  LevelInformational: 6,
  LevelDebug: 7
};

var params = {
  module:"webmonitor",    //出错的模块 应用的名称例如:xmiss
  viewurl:"",   //请求的url
  loglevel:0,  //错误等级 3err 4Warning 5Notice 7Debug
  file:"",      //出错的文件
  line:0,      //出错文件所在行
  col:0,       //出错文件所在列
  message:"",   //自定义消息
  platform:"",  //系统架构
  ua:"",              //UserAgent浏览器信息
  lang:"",           //使用的语言
  screen:"",    //分辨率
  carset:"",    //浏览器编码环境
  address:"",      //所在位置
  date:"",            //发生的时间
  timestamp:0, //发生的时间戳
}


class Log {

  constructor(params) {
    this.params = params;
  }
  //类的
  Error(message) {
    this.params.loglevel = LogLevel.LevelError;
    this.params.message = message;
  }
  Warn(message) {
    this.params.loglevel = LogLevel.LevelWarning;
    this.params.message = message;
  }
  Notice(message) {
    this.params.loglevel = LogLevel.LevelNotice;
    this.params.message = message;
  }
  //
  Info(message) {
    this.params.loglevel = LogLevel.LevelInformational;
    this.params.message = message;
  }
  Debug(message) {
    this.params.loglevel = LogLevel.LevelDebug;
    this.params.message = message;
  }
}

var log= new Log(params);

module.exports = log;
