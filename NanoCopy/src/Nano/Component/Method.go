package Component

import (
	"Nano/Session"
	"reflect"
	"unicode"
	"unicode/utf8"
)

var (
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	typeOfBytes   = reflect.TypeOf(([]byte)(nil))
	typeOfSession = reflect.TypeOf(Session.New(nil))
)

// 首字母是否为大写,是否是public
func isExported(name string) bool {
	w, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(w)
}

// isHandlerMethod decide a method is suitable handler method
func isHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}

	// Method needs three ins: receiver, *Session, []byte or pointer.
	if mt.NumIn() != 3 {
		return false
	}

	// Method needs one outs: error
	if mt.NumOut() != 1 {
		return false
	}

	if t1 := mt.In(1); t1.Kind() != reflect.Ptr || t1 != typeOfSession {
		return false
	}

	if (mt.In(2).Kind() != reflect.Ptr && mt.In(2) != typeOfBytes) || mt.Out(0) != typeOfError {
		return false
	}
	return true
}
