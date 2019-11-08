package kibanadiscover

import (
	"encoding/json"
	"testing"

	"hank.com/web-monitor/log"
)

func TestKibanadiscover(t *testing.T) {
	line := `{"module":"xmiss_group_admin","viewurl":"http://localhost:8080/#/storeOrder","filename":"storeOrderAll","enablefiledepthtype":1,"message":{"errinfo":"ReferenceError: account is not defined","bid":100022,"soid":101209,"uid":13465},"timestamp":1573205133.003}`

	//解析参数
	commonLog := &log.CommonLog{}
	err := json.Unmarshal([]byte(line), commonLog)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", commonLog)
}
