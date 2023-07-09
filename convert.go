package convert

import "syscall/js"

func ToGoType(dataValue js.Value) interface{} {
	dataType := dataValue.Type()
	var data interface{}
	switch dataType {
	case js.TypeBoolean:
		data = dataValue.Bool()
	case js.TypeNumber:
		data = dataValue.Int()
	case js.TypeObject:
		data = JSObjectToGoType(dataValue)
	case js.TypeString:
		data = dataValue.String()
	case js.TypeNull:
		data = nil
	case js.TypeUndefined:
		data = nil
	case js.TypeSymbol:
		data = dataValue.String()
	}
	return data
}

func JSObjectToGoType(value js.Value) interface{} {
	if value.InstanceOf(js.Global().Get("Uint8Array")) {
		return ToBytes(value)
	} else if js.Global().Get("Array").Get("isArray").Invoke(value).Bool() {
		return ToArray(value)
	} else if value.InstanceOf(js.Global().Get("Map")) {
		return JSMapToGoMap(value)
	}
	return JSObjectToGoMap(value)
}

func ToBytes(value js.Value) []byte {
	length := value.Length()
	bytes := make([]byte, length)
	js.CopyBytesToGo(bytes, value)
	return bytes
}

func ToArray(value js.Value) []interface{} {
	length := value.Length()
	array := make([]interface{}, length)
	for i := 0; i < length; i++ {
		array[i] = ToGoType(value.Index(i))
	}
	return array
}

func JSObjectToGoMap(value js.Value) map[string]interface{} {
	object := make(map[string]interface{})
	keys := js.Global().Get("Object").Call("keys", value)
	length := keys.Length()
	for i := 0; i < length; i++ {
		key := keys.Index(i).String()
		object[key] = ToGoType(value.Get(key))
	}
	return object
}

func JSMapToGoMap(jsMap js.Value) map[interface{}]interface{} {
	object := make(map[interface{}]interface{})
	keys := js.Global().Get("Array").Call("from", jsMap.Call("keys"))
	length := keys.Length()
	for i := 0; i < length; i++ {
		key := ToGoType(keys.Index(i))
		object[key] = ToGoType(jsMap.Call("get", key))
	}
	return object
}

func GoMapToJSObject(m map[string]interface{}) js.Value {
	obj := js.Global().Get("Object").New(nil)
	for k, v := range m {
		obj.Set(k, v)
	}
	return obj
}

type json struct {
}

func (j json) Stringify(v js.Value) string {
	return js.Global().Get("JSON").Call("stringify", v)
}

func (j json) Parse(str string) js.Value {
	return js.Global().Get("JSON").Call("parse", str)
}

var JSON = json{}
