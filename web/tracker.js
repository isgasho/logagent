const logs = require('./logs')
console.log(logs);

/**
* @des 错误对象
* @class 
*/
function E(){}

(function(window,E){
	if(!E){
		return false;
	}

	//内置参数列表
	var options = {
		bid:0,
		soid:0,
		sid:0,
		uid:0,
		url:"",
		name:"",//商户名称
		enabled:true //是否可用，关闭以后不再捕获和发送错误信息
	};

	/**
	* @des 处理时间
	* @return YYYY-MM-DD hh:mm:ss
	* @private
	*/
	function dateFun(){
		var date = new Date();
	    return date.getFullYear()+"-"+(date.getMonth() + 1)+"-"+date.getDate()+" "+date.getHours()+":"+date.getMinutes()+":"+date.getSeconds();
	}

	/**
	 * @des 获取当前时间戳
	 * @return
	 * @private
	 */
	function GetTimestamp(timeDate) {
		return parseInt(new Date(timeDate).getTime()/1000)
	}

	/**
	* @des 获取报错JS所在行的源代码上传服务器
	* @private
	* @return 报错代码
	*/
	function getCodeFun(){
		return "";
	}

	/**
	* @des 公有方法处理初始参数设置
	* @param {Error} obj, 用户传递过来的初始化信息
	* @public
	*/
	E.init = function(obj) {
		if(obj["url"] === undefined){
			return false;
		}
		for(var i in obj){
			options[i] = obj[i];
		}
	}

	/**
	* @des 公有方法处理,用于监控 `try/catch` 中被捕获的异常。
	* @param {Error} obj, 传递过来的异常对象信息。
	* @return {Object} 主要用于单元测试。
	* @public
	*/
	E.error = function(obj){
		var params_obj = {};
		if(obj instanceof Error){
			var stackArr =obj.stack.split(':');
			params_obj.info = (obj.message || obj.description) +" "+(obj.stack || obj.stacktrace),
			params_obj.module = "webmonitor";
			params_obj.line = stackArr[stackArr.length-2];//行数
			params_obj.col = stackArr[stackArr.length-1];
		}else{
			params_obj = obj;
		}
		error(params_obj);
		return true;
	}

	/**
	* @des 全局对象Errr
	* @params {String}
	* @public
	*/
	window.onerror = function(message, file, line, column,innerError){
		var params_obj = {
			info:message,
			file:file,
			line:line,
			col:column,
			module:"webmonitor"
		};
		error(params_obj);
    	return true;
	};

	//TODO
	//DOMError   DOMException

	/**
	 * @des 私有方法
	 * @param {Object} arg, 错误信息对象
	 * @private
	 */
	function error(arg){
		if(typeof arg === "string"){
			return false;
		}

		//获取时间格式
		var date =dateFun();

		var errorMsg = {
			bid:options.bid,
			soid:options.soid,
			sid:options.sid,
			uid:options.uid,
			name:options.name,//商户的ID
			code:getCodeFun(),//错误代码
			info:"无错误描述!",//错误信息
			stack:"",//堆栈错误
		};

		let message = JSON.stringify(errorMsg)

		//公共的日志格式
		var commonLogFormat={
			module:"",//模块
			errlevel:3,
			viewurl:encodeURIComponent(location.href),//URL
			file:document.currentScript.src,//出错的文件
			line:0,//出错文件所在行
			col:(window.event && window.event.errorCharacter) || 0,//出错文件所在列
			message:"",
			address:"",//用户所在位置
			platform:window.navigator.platform,//手机型号
			ua:window.navigator.userAgent.toString(),//UserAgent
			lang:navigator.language || navigator.browserLanguage || "",//使用的语言
			screen:window.screen.width+" * "+window.screen.height,//分辨率
			carset:(document.characterSet ? document.characterSet : document.charset),//浏览器编码环境
			date:date,
			timestamp:GetTimestamp(date),//发生的时间
		}

		for(var i in arg){
			if(arg[i] != ""){
				errorMsg[i] = arg[i];
			}
		}

		console.log(util.dateFun());
		return false;

		//异步上报错误
		if(options.enabled){
			setTimeout(function(){
				send({
					url:options.url,
					data:commonLogFormat
				});
			},0);
		}
	}

	/**
	 * @des 创建一个HTTP GET 请求
	 * @param {Object} obj 参数列表对象 {url:'',data:{},callback:function(){}}
	 * @private
	 */
	function send(obj){
		if(!obj.callback){
			obj.callback = function(){}
		}

		//转化json数据
		var j = JSON.stringify(obj.data)
		//base64编码
		//var data = window.btoa(j);

		url = obj.url + (obj.url.indexOf("?") < 0 ? "?" : "&") + "log="+j;
		console.log(url);

		// 忽略超长 url 请求，避免资源异常。
		if(url.length > 7713){
			return;
		}

		if(window.navigator.onLine){
			var img = new Image(1,1);
			img.onload = img.onerror = img.onabort = function(){
				obj.callback();
				img.onload = img.onerror = img.onabort = null;
				img = null;
			};
			img.src = url;
		}
	}

})(window,this.E);