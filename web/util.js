
function Util(){}

//获取当前时间
Util.prototype.dateFun = function(){
    var date = new Date();
	return date.getFullYear()+"-"+(date.getMonth() + 1)+"-"+date.getDate()+" "+date.getHours()+":"+date.getMinutes()+":"+date.getSeconds();
}

//时间转化为时间戳
Util.prototype.GetTimestamp = function(timeDate){
    return parseInt(new Date(timeDate).getTime()/1000)
}

var util = new Util();

module.exports = util;