package elastic

// ErrMonitor is a structure used for serializing/deserializing data in Elasticsearch.
type ErrMonitor struct {
	Bid       int64  `json:"bid"`  //集团Bid
	Soid      int64  `json:"soid"` //店铺Soid
	Sid       int64  `json:"sid"`
	Name      string `json:"name"`      //商户的ID
	Module    string `json:"module"`    //出错的模块
	ViewUrl   string `json:"viewUrl"`   //url
	Address   string `json:"address"`   //所在位置
	Platform  string `json:"platform"`  //系统架构
	Ua        string `json:"ua"`        //UserAgent浏览器信息
	File      string `json:"file"`      //出错的文件
	Line      int64  `json:"line"`      //出错文件所在行
	Col       int64  `json:"col"`       //出错文件所在列
	Lang      string `json:"lang"`      //使用的语言
	Screen    string `json:"screen"`    //分辨率
	Carset    string `json:"carset"`    //浏览器编码环境
	Code      string `json:"code"`      //错误代码
	Info      string `json:"info"`      //错误信息
	Date      string `json:"date"`      //发生的时间
	Timestamp int64  `json:"timestamp"` //发生的时间戳
}
