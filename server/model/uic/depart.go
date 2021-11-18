package uic

type Depart struct {
	//gorm.Model
	ID        int64   `json:"id" gorm:"primary_key;column:id"`
	Name      string  `json:"name" gorm:"column:name;type:string;size:256;index;not null;comment:'名称'"`
	Type      string  `json:"type" gorm:"column:type;type:string;size:1;comment:'机构类型'"`
	Level     string  `json:"level" gorm:"column:level;type:string;size:11;comment:'级别'"`
	Parent    *Depart `json:"parent" gorm:"-;comment:'上级机构'"`
	ParentID  int64   `json:"parentId" gorm:"column:parent_depart;comment:上级机构ID"`
	InstId    string  `json:"instId" gorm:"column:inst_id;type:string;size:10;comment:'行政机构编码'"`
	FullName  string  `json:"fullName" gorm:"column:fullname;type:string;size:1024;null;comment:'机构全称'"`
	ShortName string  `json:"shortName" gorm:"column:shortname;type:string;size:256;null;comment:'简称'"`
	Deleted   bool    `json:"deleted" gorm:"column:deleted;type:tinyint;size:1;comment:'是否删除'"`
	//FuncType string      `json:"func_type" gorm:"column:func_type;type:string;size:1;comment:'职能类型'"`
	//Order    int         `json:"order" gorm:"column:order;type:int;comment:'顺序'"`
	//ResourceNum      int       `json:"resource_num" gorm:"column:resource_num;type:int;comment:'人数'"`
	//HasChildren      bool      `json:"has_children" gorm:"column:has_children;type:tinyint;size:1;comment:'是否有子节点'"`
	//ChildrenNum      int       `json:"children_num" gorm:"column:children_num;type:int;comment:'子节点个数'"`
	//BelongKfzx       int       `json:"belong_kfzx" gorm:"column:belong_kfzx;type:int;comment:'删除标志'"`
	//FirstInstId      string    `json:"first_inst_id" gorm:"column:first_inst_id;type:string;size:10;comment:'所属一级机构编码'"`
	//TypeCode  string `json:"type_code" gorm:"column:type_code;type:string;size:9;comment:'类型编码'"`
	//Telephone string `json:"telephone" gorm:"column:telephone;type:string;size:20;comment:'电话'"`
	//Email     string `json:"email" gorm:"column:email;type:string;size:128;comment:'邮件'"`
	//SyncItdmFlag     string    `json:"sync_itdm_flag" gorm:"column:sync_itdm_flag;type:string;size:1;comment:'同步ITDM状态'"`
	//IsAdministrative bool      `json:"is_administrative" gorm:"column:is_administrative;type:tinyint;size:1;comment:'是否行政机构'"`
	//IsDepartment     bool      `json:"is_department" gorm:"column:is_department;type:tinyint;size:1;comment:'是否一级部门'"`
	//RootType         string    `json:"root_type" gorm:"column:root_type;type:string;size:2;comment:'根节点类型'"`
	//CreatedTime  time.Time `json:"created_time" gorm:"column:created_time;type:time;comment:'创建时间'"`
	//CreatedUser  int       `json:"created_user" gorm:"column:created_user;type:int;size:11;comment:'创建人'"`
	//ModifiedTime time.Time `json:"modified_time" gorm:"column:modified_time;type:time;comment:'修改时间'"`
	//ModifiedUser int       `son:"modified_user" gorm:"column:modified_user;type:int;size:11;comment:'修改人'"`
	//Abbr         string    `json:"abbr" gorm:"column:abbr;type:string;size:9;null;comment:'部门名称缩写'"`
	//CcbinsEstbDate   time.Time `json:"ccbins_estb_date" gorm:"column:ccbins_estb_date;type:time;comment:'机构创建时间'"`
	//CcbinsUndoDate   time.Time `json:"ccbins_undo_date" gorm:"column:ccbins_undo_date;type:time;comment:'机构撤销时间'"`
	//IsDismissed      string    `json:"is_dismissed" gorm:"column:is_dismissed;type:tinyint;size:1;comment:'是否解散'"`
}

func (Depart) TableName() string {
	return "kfzx"
}
