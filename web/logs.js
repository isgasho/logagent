const util = require('./util')

const LogLevel = {
  LevelEmergency: 0,
  LevelAlert: 1,
  LevelCritical: 2,
  LevelError: 3,//错误
  LevelWarning: 4,//警告
  LevelNotice: 5,
  LevelInformational: 6,//info
  LevelDebug: 7//debug
};

var levelPrefix=new Array("[M] ", "[A] ", "[C] ", "[E] ", "[W] ", "[N] ", "[I] ", "[D] ");
  
class Logs {
  constructor(params) {
    this.level = LogLevel.LevelDebug;//日志输出等级
    this.switchlog=true;  //开关

    this.module=params.module   //出错的模块 应用的名称例如:xmiss
    this.viewurl=params.viewurl
    this.loglevel=params.loglevel  //错误等级 3err 4Warning 5Notice 7Debug
    this.enablefiledepthtype=params.enablefiledepthtype
    this.filename=params.filename      //出错的文件
    this.line=params.line      //出错文件所在行
    this.col=params.col       //出错文件所在列
    this.message=params.message  //自定义消息
    this.platform=params.platform //系统架构
    this.ua=params.ua              //UserAgent浏览器信息
    this.lang=params.lang           //使用的语言
    this.screen=params.screen   //分辨率
    this.carset=params.carset    //浏览器编码环境
    this.address=params.address      //所在位置
    this.date=params.date            //发生的时间
    this.timestamp=params.timestamp //发生的时间戳
  }
  //设置输出等级
  SetLevel(l){
      this.loglevel=l;
  }
  //写入格式
  writeMsg(logLevel,msg){
      this.message = "[" + this.filename + ":" + this.line + ":" + this.col + "] " + msg;
      this.message  = levelPrefix[logLevel] + this.message;
  }
  //发送
  AsyncSend(){

  }
  Error(message) {
    if(!this.switchlog){
        return
    }
    if(LogLevel.LevelError > this.level){
        return
    }
    this.loglevel = LogLevel.LevelError;
    this.writeMsg(this.loglevel,message)
    //上报
    this.AsyncSend()
  }
  Warn(message) {
    if(!this.switchlog){
        return
    }
    if(LogLevel.LevelError > this.level){
        return
    }  
    this.loglevel = LogLevel.LevelWarning;
    this.writeMsg(this.loglevel,message)
     //上报
     this.AsyncSend()
  }
  Notice(message) {
    if(!this.switchlog){
        return
    }
    if(LogLevel.LevelError > this.level){
        return
    }
    this.loglevel = LogLevel.LevelNotice;
    this.writeMsg(this.loglevel,message)
     //上报
     this.AsyncSend()
  }
  //
  Info(message) {
    if(!this.switchlog){
        return
    }
    if(LogLevel.LevelError > this.level){
        return
    }
    this.loglevel = LogLevel.LevelInformational;
    this.writeMsg(this.loglevel,message)
     //上报
    this.AsyncSend()
  }
  Debug(message) {
    if(!this.switchlog){
        return
    }
    if(LogLevel.LevelError > this.level){
        return
    }
    this.loglevel = LogLevel.LevelDebug;
    this.writeMsg(this.loglevel,message)
     //上报
     this.AsyncSend()
  }
}

var params = {
  module:"webmonitor",    //出错的模块 应用的名称例如:xmiss
  viewurl:"",   //请求的url
  enablefiledepthtype:0,
  loglevel:0,  //错误等级 3err 4Warning 5Notice 7Debug
  filename:,      //出错的文件
  line:0,      //出错文件所在行
  col:0,       //出错文件所在列
  message:"",   //自定义消息
  platform:"",  //系统架构
  ua:"",              //UserAgent浏览器信息
  lang:"",           //使用的语言
  screen:"",    //分辨率
  carset:"",    //浏览器编码环境
  address:"",      //所在位置
  date:util.dateFun(),            //发生的时间
  timestamp:0, //发生的时间戳
}
params.timestamp=util.GetTimestamp(params.date);

var logs= new Logs(params);

module.exports = logs;

logs.Error("大事不好,出错了")

console.log(logs);

