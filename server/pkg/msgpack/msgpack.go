package msgpack

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

var encoder = msgpack.NewEncoder(nil)
var decoder = msgpack.NewDecoder(nil)

// Marshal 将 Go 对象序列化为 msgpack 格式
// 使用示例：
//
//	type User struct {
//	    ID   int    `msgpack:"id"`
//	    Name string `msgpack:"name"`
//	}
//	user := User{ID: 1, Name: "Alice"}
//	data, err := msgpack.Marshal(user)
func Marshal(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Unmarshal 将 msgpack 字节数组反序列化为 Go 对象
// 使用示例：
//
//	type User struct {
//	    ID   int    `msgpack:"id"`
//	    Name string `msgpack:"name"`
//	}
//	var user User
//	err := msgpack.Unmarshal(data, &user)
func Unmarshal(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}

// UnmarshalToJson 将 msgpack 数据转换为可读的字符串（JSON 格式）
// 使用示例：
//
//	msgpackData := []byte{...} // msgpack 格式数据
//	str, err := msgpack.UnmarshalToString(msgpackData)
//	fmt.Println(str) // {"id":1,"name":"Alice"}
func UnmarshalToJson(data []byte) (string, error) {
	var v any
	if err := msgpack.Unmarshal(data, &v); err != nil {
		return "", err
	}
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// JSONToMsgpack 将 JSON 数据转换为 msgpack 格式
// 使用示例：
//
//	jsonData := []byte(`{"id":1,"name":"Alice"}`)
//	msgpackData, err := msgpack.JSONToMsgpack(jsonData)
func JSONToMsgpack(data []byte) ([]byte, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return msgpack.Marshal(v)
}
