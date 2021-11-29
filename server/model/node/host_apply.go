package node

import "github.com/alexreagan/rabbit/server/model/gtime"

const HostApplyStateSubmitted = "submitted"
const HostApplyStateSuccess = "success"
const HostApplyStateFailure = "failure"

type HostApplyRequest struct {
	ID        int64       `json:"id" gorm:"primary_key;column:id"`
	Name      string      `json:"name" gorm:"column:name;type:string;size:128;comment:"`
	CPU       int64       `json:"cpu" gorm:"column:cpu;comment:CPU需求数（核）"`
	Memory    int64       `json:"memory" gorm:"column:memory;comment:Memory需求数（GB）"`
	Count     int64       `json:"count" gorm:"column:count;comment:机器需求数量"`
	Remark    string      `json:"remark" gorm:"column:remark;type:string;size:512;comment:备注"`
	Creator   string      `json:"creator" gorm:"column:creator;type:string;size:128;comment:创建人"`
	Applier   string      `json:"applier" gorm:"column:applier;type:string;size:128;comment:申请人"`
	State     string      `json:"state" gorm:"column:state;type:enum('submitted', 'failure','success');default:submitted;comment:机器状态"`
	CreateAt  gtime.GTime `json:"createAt" gorm:"column:create_at;default:null;comment:创建时间"`
	ReleaseAt gtime.GTime `json:"releaseAt" gorm:"column:release_at;default:null;comment:归还时间"`
}

func (r HostApplyRequest) TableName() string {
	return "host_apply_request"
}
