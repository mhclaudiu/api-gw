package functions

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"
)

func ExecTime(t time.Time) string {

	//buff := fmt.Sprintf("%v", time.Now().Sub(t))
	buff := fmt.Sprintf("%v", time.Since(t))

	//re1 := regexp.MustCompile("^(?<=).+?(?=\\.)")
	//re2 := regexp.MustCompile("[^.]+(?=[a-zA-Z])")
	//re3 := regexp.MustCompile("[a-zA-Z]+(?=)$")
	re1 := regexp.MustCompile(`^[0-9]+`)
	re2 := regexp.MustCompile(`\.[\d]+`)
	re3 := regexp.MustCompile(`[\x41-\xff]+$`)

	f1 := re1.FindString(buff)
	f2 := re2.FindString(buff)
	f3 := re3.FindString(buff)

	//log.Printf("%s | %s-%s-%s", buff, f1, f2, f3)

	//f2 = strings.Replace(f2, ".", "", -1)

	if len(f2) > 3 {
		f2 = f2[0:4]
	}

	if len(f1) < 2 {
		f1 = fmt.Sprintf("0%s", f1)
	}

	return fmt.Sprintf("%s%s%s", f1, f2, f3)

	/*l := len(buff)
	if l > 10 {
		if buff[l-2] == 'm' {
			buff = buff[0:l-5] + "ms"
		} else {
			buff = buff[0:l-7] + "s"
		}
	}*/
}

func Split(s, d string) ([]string, error) {
	r := []string{"N/A Split"}
	if len(s) < 3 {
		return r, errors.New(r[0])
	}
	if d == "." && (s[0] >= '0' && s[0] <= '9') {
		return []string{s}, nil
	}
	r = strings.Split(s, d)
	if len(r) > 0 {
		return r, nil
	}
	return r, errors.New(r[0])
}

func RemoteHost(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if len(ip) < 4 {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if len(ip) < 4 {
		ip = r.RemoteAddr
	}
	s, _ := Split(ip, ".")
	return s[0]
}

func RemoteIP(addr string) string {

	s, _ := Split(addr, ":")
	return s[0]
}

func Dump(d interface{}) string {

	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	return string(b)
}

func GetSeconds(period string) (error, int) {

	period = strings.Trim(period, " ")

	if len(period) < 1 || period[0] == '0' {

		return nil, 0
	}

	buff := period[:len(period)-1]

	atoi, err := strconv.Atoi(buff)

	if err != nil {

		return err, 0
	}

	switch period[len(period)-1:] {

	case "m":

		return nil, atoi * 60

	case "h":

		return nil, atoi * 3600

	}

	return nil, atoi
}

func GetMiliSeconds(period string) (error, float64) {

	period = strings.Trim(period, " ")

	if len(period) < 1 || period[0] == '0' {

		return nil, 0
	}

	buff := period[:len(period)-1]

	ftoi, err := strconv.ParseFloat(buff, 64)

	if err != nil {

		return err, 0
	}

	ftoi *= 1000

	switch period[len(period)-1:] {

	case "m":

		return nil, ftoi * 60

	case "h":

		return nil, ftoi * 3600

	}

	return nil, ftoi
}

func HashSHA256(data string) string {

	h := sha256.New()
	h.Write([]byte(data))

	//aa := sha256.Sum256([]byte("pass"))
	//log.Printf("%x", aa[:])

	return hex.EncodeToString(h.Sum(nil))
}

func HashSHA1(data string) string {

	h := sha1.New()
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

func escape(source string) string {
	var j int = 0
	if len(source) == 0 {
		return ""
	}
	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}

func Escape(sql string) string {

	sql = strings.TrimSpace(sql)
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]

		escape = 0

		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
			break
		case '\n': /* Must be escaped for logs */
			escape = 'n'
			break
		case '\r':
			escape = 'r'
			break
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		case '"': /* Better safe than sorry */
			escape = '"'
			break
		case '\032': //十进制26,八进制32,十六进制1a, /* This gives problems on Win32 */
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}

func MysqlRealEscapeString(value string) string {
	var sb strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch c {
		case '\\', 0, '\n', '\r', '\'', '"':
			sb.WriteByte('\\')
			sb.WriteByte(c)
		case '\032':
			sb.WriteByte('\\')
			sb.WriteByte('Z')
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

func ToMap(in interface{}) string {
	j, _ := json.Marshal(in)
	return string(j)
}

func GenerateOTP(maxDigits int) string {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%0*d", maxDigits, bi)
}

func generateRandomNumber(numberOfDigits int) (int, error) {
	maxLimit := int64(int(math.Pow10(numberOfDigits)) - 1)
	lowLimit := int(math.Pow10(numberOfDigits - 1))

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(maxLimit))
	if err != nil {
		return 0, err
	}
	randomNumberInt := int(randomNumber.Int64())

	// Handling integers between 0, 10^(n-1) .. for n=4, handling cases between (0, 999)
	if randomNumberInt <= lowLimit {
		randomNumberInt += lowLimit
	}

	// Never likely to occur, kust for safe side.
	if randomNumberInt > int(maxLimit) {
		randomNumberInt = int(maxLimit)
	}
	return randomNumberInt, nil
}

func GenerateOTP2(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpChars := "1234567890"

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func GenCaptchaCode() (string, error) {
	codes := make([]byte, 6)
	if _, err := rand.Read(codes); err != nil {
		return "", err
	}

	for i := 0; i < 6; i++ {
		codes[i] = uint8(48 + (codes[i] % 10))
	}

	return string(codes), nil
}

func GenerateCode() string {
	return fmt.Sprint(time.Now().Nanosecond())[:6]
}

func GenerateOTPCode(length int) (string, error) {
	seed := "012345679"
	byteSlice := make([]byte, length)

	for i := 0; i < length; i++ {
		max := big.NewInt(int64(len(seed)))
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		byteSlice[i] = seed[num.Int64()]
	}

	return string(byteSlice), nil
}

func FormatRemaining(total int) string {

	/*h := total / 3600
	s := total % 60
	m := (total / 60) - h*60*/

	s := total % 60
	total /= 60
	m := total % 60
	total /= 60

	return fmt.Sprintf("%02dh:%02dm:%02ds", total, m, s)
}

func FormatUptime(total int) string {

	days := total / (60 * 60 * 24)
	hours := (total - (days * 60 * 60 * 24)) / (60 * 60)
	minutes := ((total - (days * 60 * 60 * 24)) - (hours * 60 * 60)) / 60
	seconds := total % 60

	return fmt.Sprintf("%02dd:%02dh:%02dm:%02ds", days, hours, minutes, seconds)
}

func TimerWait(duration int) {

	colored := color.New(color.FgHiBlue, color.Italic).SprintfFunc()

	ticker := time.Tick(time.Second)

	ticker2 := time.Tick(time.Millisecond * 10)
	s := ""
	for i := duration; i >= 0; i-- {

		<-ticker

		WG.Add(1)

		go func() {

			for j := 100; j >= 0; j-- {

				<-ticker2

				/*switch {

				case j < 100:

					s = fmt.Sprintf("0%d", j)

				case j < 10:

					s = fmt.Sprintf("00%d", j)
				}*/

				if j < 100 {

					s = fmt.Sprintf("0%d", j)
				}

				if j < 10 {

					s = fmt.Sprintf("00%d", j)
				}

				if i == 0 {
					j = 0
				}

				fmt.Printf("\r --> %s", colored("0%ds:%sms", i, s))
			}

			WG.Done()
		}()

		WG.Wait()
	}
	fmt.Println()
}

func PointerTo[T any](v T) *T {
	return &v
}

func GetMaxString(global, local string) string {

	if len(local) > 1 && local[0] != '0' {

		return local
	}

	return global
}

func CheckPointerBool(b *bool) bool {

	if b != nil {

		return *b
	}

	return false
}

func CheckPointerInt(i *int) int {

	if i != nil {

		return *i
	}

	return -1
}

func FirstToUpper(s string) string {

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func ValidateRegexFormat(s, f string) bool {

	re, err := regexp.Compile(f)

	if err != nil {

		log.Panic(err)
	}

	if !re.MatchString(s) {

		return false
	}

	return true
}

// 55.693334ms
// ^[\d]+\.[\d]{1,3}|[a-z]+$
func ParsePingMs(s string) (error, string) {

	re, err := regexp.Compile(`^[\d]+\.[\d]{1,3}|^[\d]+`)
	if err != nil {

		return err, s
	}

	/*re2, err := regexp.Compile(`[a-zA-Z]+$`)
	if err != nil {

		return err, s
	}

	if re.MatchString(s) && re2.MatchString(s) {

		return nil, re.FindString(s) + re2.FindString(s)
	}*/

	if re.MatchString(s) {

		return nil, re.FindString(s)
	}

	return nil, s
}

func ParseInterval(period string) (error, string) {

	period = strings.Trim(period, " ")

	atoi, err := strconv.Atoi(period[:len(period)-1])

	if err != nil {

		return err, ""
	}

	switch period[len(period)-1:] {

	case "m":

		return nil, fmt.Sprintf("%d months", atoi)

	}

	return nil, fmt.Sprintf("%d days", atoi)
}
