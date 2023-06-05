package vingo

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"unicode"
)

func TreeBuildString(list *[]map[string]any, id string, pidName string) (result []map[string]any) {
	for _, row := range *list {

		if row[pidName] != id {
			continue
		}

		children := TreeBuildString(list, row["id"].(string), pidName)

		if len(children) > 0 {
			row["hasChild"] = true
			row["children"] = children
		} else {
			row["hasChild"] = false
		}
		row["id"] = GetUUID()
		result = append(result, row)
	}
	return
}

func TreeBuild(list *[]map[string]any, id uint, pidName string, already *[]uint) (result []map[string]any) {

	for _, row := range *list {

		if ToUint(row[pidName]) != id {
			continue
		}

		*already = append(*already, id)

		children := TreeBuild(list, ToUint(row["id"]), pidName, already)

		if len(children) > 0 {
			row["hasChild"] = true
			row["children"] = children
		} else {
			row["hasChild"] = false
		}
		result = append(result, row)
	}
	return
}

func TreeBuilds(list *[]map[string]any, ids []uint, pidName string) []map[string]any {
	result := make([]map[string]any, 0)
	already := make([]uint, 0)
	for _, id := range ids {
		if IsInSlice(id, already) {
			continue
		}
		result = append(result, TreeBuild(list, id, pidName, &already)...)
	}
	return result
}

func CallStructFunc(obj any, method string, param map[string]any) any {
	t := reflect.TypeOf(obj)
	_func, ok := t.MethodByName(method)
	if !ok {
		panic(fmt.Sprintf("%v方法不存在", method))
	}

	_param := make([]reflect.Value, 0)
	_param = append(_param, reflect.ValueOf(obj))
	for _, value := range param {
		_param = append(_param, reflect.ValueOf(value))
	}
	res := _func.Func.Call(_param)
	return res[0].Interface()
}

func CheckErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func SqlLike(keyword string) string {
	return fmt.Sprintf("%%%v%%", strings.Trim(keyword, " "))
}

/**
 * 增长率计算
 * @param float64 $now 现在值
 * @param float64 $prev 过去值
 * @return string
 */
func ComputeGrowRate(now float64, prev float64) string {
	if now == prev {
		return "0.00"
	} else if prev == 0 {
		return "-"
	} else {
		return fmt.Sprintf("%.2f", ((now - prev) / prev * 100))
	}
}

/**
 * 根据起点坐标和终点坐标测距离
 * @param [2]float64 $from [起点坐标(经纬度),例如:[2]float64{118.012951,36.810024}]
 * @param [2]float64 $to [终点坐标(经纬度)]
 * @param bool $km 是否以公里为单位 false:米 true:公里(千米)
 * @param int $decimal 精度 保留小数位数
 * @return float  距离数值
 */
func Distance(from Location, to Location, km bool, decimal int) float64 {
	EARTH_RADIUS := 6370.996 // 地球半径系数
	fromSorted := from
	toSorted := to
	if from.Lng > to.Lng {
		fromSorted = to
		toSorted = from
	}

	dLat := (toSorted.Lng - fromSorted.Lng) * math.Pi / 180
	dLon := (toSorted.Lat - fromSorted.Lat) * math.Pi / 180

	fromLat := fromSorted.Lng * math.Pi / 180
	toLat := toSorted.Lng * math.Pi / 180

	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(fromLat)*math.Cos(toLat)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Asin(math.Sqrt(a))

	distance := EARTH_RADIUS * c * 1000

	if km {
		distance = distance / 1000
	}

	return math.Round(distance*math.Pow10(decimal)) / math.Pow10(decimal)
}

// 密码加密
func PasswordToCipher(text string, salt string) string {
	return MD5(MD5(text) + salt)
}

// 密码强度验证
// level-2：中等密码，任意两种字符组合
// level-3：复杂密码，必须包含四种字符组合
func PasswordStrength(password string, level int) {
	if len(password) < 6 || len(password) > 18 {
		panic("密码长度需符合6-18个字符长度要求")
	}
	if level == 2 {
		// 中等密码，任意两种字符组合
		hasDigit := false
		hasUpper := false
		hasLower := false
		hasSpecial := false

		for _, ch := range password {
			if unicode.IsDigit(ch) {
				hasDigit = true
			} else if unicode.IsUpper(ch) {
				hasUpper = true
			} else if unicode.IsLower(ch) {
				hasLower = true
			} else if unicode.IsPunct(ch) || unicode.IsSymbol(ch) {
				hasSpecial = true
			}
		}
		if !(hasDigit && hasUpper) && !(hasDigit && hasLower) && !(hasDigit && hasSpecial) &&
			!(hasUpper && hasLower) && !(hasUpper && hasSpecial) && !(hasLower && hasSpecial) {
			panic("密码需满足两种以上的字符组合（数字、大写字母、小写字母、特殊符号）")
		}
	} else if level == 3 {
		// 复杂密码，必须包含四种字符组合
		hasDigit := false
		hasUpper := false
		hasLower := false
		hasSpecial := false

		for _, ch := range password {
			if unicode.IsDigit(ch) {
				hasDigit = true
			} else if unicode.IsUpper(ch) {
				hasUpper = true
			} else if unicode.IsLower(ch) {
				hasLower = true
			} else if unicode.IsPunct(ch) || unicode.IsSymbol(ch) {
				hasSpecial = true
			}
		}
		if !hasDigit || !hasUpper || !hasLower || !hasSpecial {
			panic("密码需满足四种字符组合（数字、大写字母、小写字母、特殊符号）")
		}
	}
}

// 返回传入参数的指针
func StringPointer(text string) *string {
	return &text
}
