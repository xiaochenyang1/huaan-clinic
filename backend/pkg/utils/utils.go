package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateShortUUID 生成短UUID（去掉横线）
func GenerateShortUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// GenerateAppointmentNo 生成预约编号
// 格式：YYYYMMDDHHmmss + 6位随机数
func GenerateAppointmentNo() string {
	now := time.Now()
	dateStr := now.Format("20060102150405")
	randomStr := GenerateRandomNumber(6)
	return dateStr + randomStr
}

// GenerateRandomNumber 生成指定长度的随机数字字符串
func GenerateRandomNumber(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	for i := 0; i < length; i++ {
		result[i] = digits[randomBytes[i]%10]
	}
	return string(result)
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	for i := 0; i < length; i++ {
		result[i] = charset[randomBytes[i]%byte(len(charset))]
	}
	return string(result)
}

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// MD5 计算MD5哈希
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// ValidatePhone 验证手机号格式
func ValidatePhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// ValidateUsername 验证用户名格式
// 规则：4-20字符，只能包含字母、数字、下划线，必须以字母开头
func ValidateUsername(username string) bool {
	pattern := `^[a-zA-Z][a-zA-Z0-9_]{3,19}$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

// ValidateIDCard 验证身份证号格式（简单校验）
func ValidateIDCard(idCard string) bool {
	// 18位身份证号
	pattern := `^\d{17}[\dXx]$`
	matched, _ := regexp.MatchString(pattern, idCard)
	return matched
}

// MaskPhone 手机号脱敏（显示前3后4）
func MaskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

// MaskIDCard 身份证号脱敏（显示前6后4）
func MaskIDCard(idCard string) string {
	if len(idCard) < 10 {
		return idCard
	}
	return idCard[:6] + "********" + idCard[len(idCard)-4:]
}

// MaskName 姓名脱敏
func MaskName(name string) string {
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

// ParseDate 解析日期字符串（YYYY-MM-DD）
func ParseDate(dateStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", dateStr, time.Local)
}

// FormatDate 格式化日期（YYYY-MM-DD）
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime 格式化日期时间（YYYY-MM-DD HH:mm:ss）
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatTime 格式化时间（HH:mm）
func FormatTime(t time.Time) string {
	return t.Format("15:04")
}

// GetTodayStart 获取今天开始时间
func GetTodayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

// GetTodayEnd 获取今天结束时间
func GetTodayEnd() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.Local)
}

// GetMonthStart 获取本月开始时间
func GetMonthStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
}

// GetMonthEnd 获取本月结束时间
func GetMonthEnd() time.Time {
	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)
	return time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local).Add(-time.Nanosecond)
}

// CalculateAge 根据身份证号计算年龄
func CalculateAge(idCard string) int {
	if len(idCard) != 18 {
		return 0
	}
	birthYear, _ := strconv.Atoi(idCard[6:10])
	birthMonth, _ := strconv.Atoi(idCard[10:12])
	birthDay, _ := strconv.Atoi(idCard[12:14])

	now := time.Now()
	age := now.Year() - birthYear

	// 如果还没过生日，年龄减1
	if int(now.Month()) < birthMonth || (int(now.Month()) == birthMonth && now.Day() < birthDay) {
		age--
	}

	return age
}

// GetGenderFromIDCard 从身份证号获取性别
func GetGenderFromIDCard(idCard string) int {
	if len(idCard) != 18 {
		return 0
	}
	genderCode, _ := strconv.Atoi(string(idCard[16]))
	if genderCode%2 == 0 {
		return 2 // 女
	}
	return 1 // 男
}

// GetBirthDateFromIDCard 从身份证号获取出生日期
func GetBirthDateFromIDCard(idCard string) string {
	if len(idCard) != 18 {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", idCard[6:10], idCard[10:12], idCard[12:14])
}

// ContainsChinese 检查字符串是否包含中文
func ContainsChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

// TrimAllSpaces 去除所有空格
func TrimAllSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

// InArray 检查元素是否在切片中
func InArray[T comparable](item T, array []T) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}
