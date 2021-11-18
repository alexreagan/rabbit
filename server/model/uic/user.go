package uic

import (
	"github.com/spf13/viper"
)

type User struct {
	//gorm.Model
	ID          int64  `json:"id" gorm:"primary_key;column:id"`
	UserName    string `json:"username" gorm:"column:username;type:string;size:80;unique;not null;comment:用户名"`
	CnName      string `json:"cnName" gorm:"column:nickname;type:string;size:80;null;comment:中文名"`
	Password    string `json:"password" gorm:"column:password;type:string;size:128;not null;comment:密码"`
	Sex         string `json:"sex" gorm:"column:sex;type:string;size:1;null;comment:性别"`
	Mobile      string `json:"mobile" gorm:"column:mobile;type:string;size:81;index;null;comment:手机"`
	Telephone   string `json:"telephone" gorm:"column:telphone;type:string;size:80;null;comment:电话"`
	Email       string `json:"email" gorm:"column:email;type:string;size:1024;null;comment:邮箱"`
	InstId      string `json:"instId" gorm:"column:inst_id;type:string;size:9;null;comment:所属机构编码"`
	FirstInstId string `json:"firstInstId" gorm:"column:first_inst_id;type:string;size:9;null;comment:所属一级机构编码"`
	JgygUserId  string `json:"jgygUserId" gorm:"column:jgyg_user_id;type:string;size:8;unique;index;no null;comment:机构员工编号"`
	AdUserName  string `json:"adUserName" gorm:"column:ad_username;type:string;size:80;unique;index;not null;comment:AD(云桌面)用户名"`
	IsSuperuser bool   `json:"isSuperuser" gorm:"column:is_superuser;type:tinyint;size:1;default:0;comment:是否为超级用户"`
	//Birthday    time.Time `json:"birthday" gorm:"column:birthday;type:time;null;comment:出生日期"`
	//Image       []byte    `json:"image" gorm:"column:image;type:bytes;null;comment:头像"`
	//LastLogin   time.Time `json:"lastLogin" gorm:"column:last_login;type:time;null;comment:上次登录时间"`
	//Age         uint8     `gorm:"column:age;type:uint;null;comment:年龄"`
	//IdentityCard string                  `gorm:"column:identity_card;type:string;size:18;null;comment:身份证号"`
	//Departments  []department.Depart `gorm:"many2many:user_department_rel;"`
	//Kfzx             string    `gorm:"type:varchar(10);null"`
	//BelongDepartment string    `form:"belong_department" gorm:"type:varchar(20);null"`
	//Org              string    `form:"org" json:"org" gorm:"type:varchar(20);null"`
	//Tx               string    `form:"tx" json:"tx" gorm:"type:varchar(20);null"`
	//IsManager        bool      `form:"is_manager" json:"is_manager" gorm:"type:tinyint(1);default:0;comment:"`
	//IsGM             bool      `form:"is_gm" json:"is_gm" gorm:"type:tinyint(1);default:0;comment:"`
	//Manager          int       `form:"manager" json:"manager" gorm:"type:int(11);null"`
	//Post             string    `form:"post" json:"post" gorm:"type:varchar(40);null;comment:岗位"`
	//Duty             string    `form:"duty" json:"duty" gorm:"type:varchar(40);null;comment:职务"`
	//KqId string `json:"kq_id" gorm:"column:kfid;type:string;size:40;null;comment:考勤编号"`
	//Type             string    `gorm:"type:varchar(2);null;comment:人员属性"`
	//ExpiryDate       time.Time `json:"expiry_date" gorm:"type:datetime;null;comment:失效期"`
	//Company          string    `json:"company" gorm:"type:int(11);null;comment:所属公司机构"`
	//AccessLevel      string    `json:"access_level" gorm:"type:varchar(1);null;comment:权限级别"`
	//Order            int       `json:"order" gorm:"type:int(11);default:9000;null;comment:顺序号"`
	//Memo string `json:"memo" gorm:"column:memo;type:string;size:256;null;comment:备注"`
	//IsBiz            string    `json:"is_biz" gorm:"type:varchar(1);default:'T';null;comment:技术业务标志"`
	//IsSyncAd         string    `json:"is_sync_ad" gorm:"type:varchar(1);null;comment:是否同步到AD"`
	//CreatedTime  time.Time `json:"created_time" gorm:"column:created_time;type:time;null;comment:创建时间"`
	//CreatedUser  int       `json:"created_user" gorm:"column:created_user;type:int;size:11;null;comment:创建人"`
	//ModifiedTime time.Time `json:"modified_time" gorm:"column:modified_time;type:time;null;comment:修改时间"`
	//ModifiedUser int       `json:"modified_user" gorm:"column:modified_user;type:int;size:11;null;comment:修改人"`
	//DeletedTime  time.Time `json:"deleted_time" gorm:"column:deleted_time;type:time;null;comment:删除时间"`
	//DefaultOrg   int    `json:"default_org" gorm:"type:int(11);null;comment:默认部门"`
	//JgygUserName string `json:"jgyg_username" gorm:"column:jgyg_username;type:string;size:100;null;comment:机构员工用户名"`
	//HrUserId     string `json:"hr_user_id" gorm:"column:hr_user_id;type:string;size:15;null;comment:人力资源员工编号"`
	//IsJgyg       bool   `json:"is_jgyg" gorm:"type:tinyint(1);null;comment:是否机构员工用户"`
	//InstSn      int    `json:"inst_sn" gorm:"column:inst_sn;type:int;null;comment:所属二级部门ID"`
}

func (this User) TableName() string {
	return "user"
}

func skipAccessControl() bool {
	return !viper.GetBool("access_control")
}

func (this User) IsAdmin() bool {
	if skipAccessControl() {
		return true
	}
	return this.IsSuperuser
}

type Session struct {
	ID      int64
	Uid     int64
	Sig     string
	Expired int
}

func (this Session) TableName() string {
	return "session"
}
