package model

import (
	"time"
)

// MedicalRecord 就诊记录模型
type MedicalRecord struct {
	BaseModel
	AppointmentID int64     `gorm:"uniqueIndex;not null;comment:预约ID" json:"appointment_id"`
	PatientID     int64     `gorm:"index;not null;comment:患者ID" json:"patient_id"`
	DoctorID      int64     `gorm:"index;not null;comment:医生ID" json:"doctor_id"`
	DepartmentID  int64     `gorm:"index;not null;comment:科室ID" json:"department_id"`
	VisitDate     time.Time `gorm:"type:date;index;not null;comment:就诊日期" json:"visit_date"`
	Diagnosis     string    `gorm:"type:text;comment:诊断结果" json:"diagnosis"`
	Prescription  string    `gorm:"type:text;comment:处方" json:"prescription"`
	Advice        string    `gorm:"type:text;comment:医嘱" json:"advice"`
	Remark        string    `gorm:"type:varchar(512);comment:备注" json:"remark"`

	// 关联
	Appointment *Appointment `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`
	Patient     *Patient     `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Doctor      *Doctor      `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	Department  *Department  `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

// TableName 表名
func (MedicalRecord) TableName() string {
	return "medical_records"
}

// MedicalRecordVO 就诊记录视图对象
type MedicalRecordVO struct {
	ID             int64  `json:"id"`
	AppointmentID  int64  `json:"appointment_id"`
	AppointmentNo  string `json:"appointment_no,omitempty"`
	PatientID      int64  `json:"patient_id"`
	PatientName    string `json:"patient_name"`
	DoctorID       int64  `json:"doctor_id"`
	DoctorName     string `json:"doctor_name"`
	DoctorTitle    string `json:"doctor_title"`
	DepartmentID   int64  `json:"department_id"`
	DepartmentName string `json:"department_name"`
	VisitDate      string `json:"visit_date"`
	Diagnosis      string `json:"diagnosis"`
	Prescription   string `json:"prescription"`
	Advice         string `json:"advice"`
	Remark         string `json:"remark,omitempty"`
	CreatedAt      string `json:"created_at"`
}

// ToVO 转换为视图对象
func (r *MedicalRecord) ToVO() *MedicalRecordVO {
	vo := &MedicalRecordVO{
		ID:            r.ID,
		AppointmentID: r.AppointmentID,
		PatientID:     r.PatientID,
		DoctorID:      r.DoctorID,
		DepartmentID:  r.DepartmentID,
		VisitDate:     r.VisitDate.Format("2006-01-02"),
		Diagnosis:     r.Diagnosis,
		Prescription:  r.Prescription,
		Advice:        r.Advice,
		Remark:        r.Remark,
		CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if r.Appointment != nil {
		vo.AppointmentNo = r.Appointment.AppointmentNo
	}
	if r.Patient != nil {
		vo.PatientName = maskName(r.Patient.Name)
	}
	if r.Doctor != nil {
		vo.DoctorName = r.Doctor.Name
		vo.DoctorTitle = GetTitleName(r.Doctor.Title)
	}
	if r.Department != nil {
		vo.DepartmentName = r.Department.Name
	}

	return vo
}

// MedicalRecordListVO 就诊记录列表视图对象
type MedicalRecordListVO struct {
	ID             int64  `json:"id"`
	PatientName    string `json:"patient_name"`
	DoctorName     string `json:"doctor_name"`
	DepartmentName string `json:"department_name"`
	VisitDate      string `json:"visit_date"`
	Diagnosis      string `json:"diagnosis"`
}

// ToListVO 转换为列表视图对象
func (r *MedicalRecord) ToListVO() *MedicalRecordListVO {
	vo := &MedicalRecordListVO{
		ID:        r.ID,
		VisitDate: r.VisitDate.Format("2006-01-02"),
		Diagnosis: r.Diagnosis,
	}

	if r.Patient != nil {
		vo.PatientName = maskName(r.Patient.Name)
	}
	if r.Doctor != nil {
		vo.DoctorName = r.Doctor.Name
	}
	if r.Department != nil {
		vo.DepartmentName = r.Department.Name
	}

	return vo
}
