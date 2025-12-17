package model

// Patient 就诊人模型
type Patient struct {
	BaseModel
	UserID    int64  `gorm:"index;not null;comment:所属用户ID" json:"user_id"`
	Name      string `gorm:"type:varchar(32);not null;comment:姓名" json:"name"`
	IDCard    string `gorm:"type:varchar(18);index;not null;comment:身份证号" json:"id_card"`
	Phone     string `gorm:"type:varchar(20);not null;comment:手机号" json:"phone"`
	Gender    int    `gorm:"type:tinyint;default:0;comment:性别 0未知 1男 2女" json:"gender"`
	BirthDate string `gorm:"type:date;comment:出生日期" json:"birth_date"`
	Relation  string `gorm:"type:varchar(20);default:'self';comment:与用户关系" json:"relation"`
	IsDefault int    `gorm:"type:tinyint;default:0;comment:是否默认就诊人 0否 1是" json:"is_default"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 表名
func (Patient) TableName() string {
	return "patients"
}

// PatientVO 就诊人视图对象
type PatientVO struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	IDCard       string `json:"id_card"` // 脱敏后的身份证号
	Phone        string `json:"phone"`   // 脱敏后的手机号
	Gender       int    `json:"gender"`
	GenderName   string `json:"gender_name"`
	BirthDate    string `json:"birth_date"`
	Age          int    `json:"age"`
	Relation     string `json:"relation"`
	RelationName string `json:"relation_name"`
	IsDefault    int    `json:"is_default"`
	LastLoginIP  string `json:"last_login_ip,omitempty"` // 用户最后登录IP
}

// ToVO 转换为视图对象
func (p *Patient) ToVO() *PatientVO {
	return &PatientVO{
		ID:           p.ID,
		Name:         maskName(p.Name),
		IDCard:       maskIDCard(p.IDCard),
		Phone:        maskPhone(p.Phone),
		Gender:       p.Gender,
		GenderName:   getGenderName(p.Gender),
		BirthDate:    p.BirthDate,
		Age:          calculateAge(p.IDCard),
		Relation:     p.Relation,
		RelationName: GetRelationName(p.Relation),
		IsDefault:    p.IsDefault,
	}
}

// ToFullVO 转换为完整视图对象（不脱敏，用于编辑）
func (p *Patient) ToFullVO() *PatientVO {
	return &PatientVO{
		ID:           p.ID,
		Name:         p.Name,
		IDCard:       p.IDCard,
		Phone:        p.Phone,
		Gender:       p.Gender,
		GenderName:   getGenderName(p.Gender),
		BirthDate:    p.BirthDate,
		Age:          calculateAge(p.IDCard),
		Relation:     p.Relation,
		RelationName: GetRelationName(p.Relation),
		IsDefault:    p.IsDefault,
	}
}

// maskName 姓名脱敏
func maskName(name string) string {
	runes := []rune(name)
	if len(runes) <= 1 {
		return name
	}
	if len(runes) == 2 {
		return string(runes[0]) + "*"
	}
	masked := string(runes[0])
	for i := 1; i < len(runes)-1; i++ {
		masked += "*"
	}
	masked += string(runes[len(runes)-1])
	return masked
}

// maskIDCard 身份证脱敏
func maskIDCard(idCard string) string {
	if len(idCard) < 10 {
		return idCard
	}
	return idCard[:6] + "********" + idCard[len(idCard)-4:]
}

// getGenderName 获取性别名称
func getGenderName(gender int) string {
	switch gender {
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	default:
		return "未知"
	}
}

// calculateAge 根据身份证计算年龄
func calculateAge(idCard string) int {
	if len(idCard) != 18 {
		return 0
	}
	// 简单实现，实际应该用time包计算
	birthYear := 0
	for i := 6; i < 10; i++ {
		birthYear = birthYear*10 + int(idCard[i]-'0')
	}
	// 假设当前年份，实际应该用time.Now().Year()
	currentYear := 2025
	return currentYear - birthYear
}
