package logging

import (
	"api-gw/functions"
	"fmt"
	"regexp"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GetSeverity(code int) int {

	if code >= 100 {

		return code
	}

	if code == 0 || code > 1 {

		return CONST_CODE_WARNING
	}

	if code < 0 {

		return CONST_CODE_ERROR
	}

	return CONST_CODE_INFO
}

func (l Log) ProcessEventID(s *string, code int) error {

	re, err := regexp.Compile(`error|Err|Error|SQL\sErr|SQL\sError|#mask`)

	if err != nil {

		l.Add(Entry{
			Err: err,
		})

		return nil
	}

	if re.MatchString(*s) {

		re2, err := regexp.Compile(`\s?#mask\s?$`)

		if err != nil {

			l.Add(Entry{
				Err: err,
			})

			return nil
		}

		hash := functions.HashSHA256(*s)

		//*s = strings.Replace(*s, " #mask", "", 1)
		//int(math.Abs(float64(code)))
		if re2.MatchString(*s) {

			*s = re2.ReplaceAllString(*s, "")
		}

		hashLen := CONST_ERROR_ID_HASH_LEN

		if functions.CheckPointerInt(l.ErrorIDLength) > 0 {
			hashLen = *l.ErrorIDLength
		}

		//return fmt.Errorf("0%dx%s", code, cases.Title(language.English, cases.NoLower).String(uniuri.NewLen(20)))
		return fmt.Errorf("0%dx%s", code, cases.Title(language.English, cases.NoLower).String(hash[:hashLen]))
	}

	return nil
}
