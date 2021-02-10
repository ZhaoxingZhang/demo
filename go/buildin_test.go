package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)
type sheet struct {
	Token  string   `json:"token"`
	Blocks []string `json:"blocks"`
}
const (
	JsonFilePath = "data/data.json"
	FakeJsonFilePath = "data/fake.json"
)
func ParseJson(filePath string) {
	var err error
	var st sheet
	pwd, err := os.Getwd()
	data, err := ioutil.ReadFile(filepath.Join(pwd,filePath))
	err = json.Unmarshal(data, &st)
	if err != nil {
		errMsg := fmt.Sprintf("failed to loadUserDataFromFile: %v, error: %v", filePath, err)
		fmt.Errorf(errMsg)
	}
	fmt.Printf("file %v, st %v, more is ignored\n", filePath, st)
}
func TestJson(t *testing.T) {
	ParseJson(JsonFilePath)
	ParseJson(FakeJsonFilePath)
}
// strings.Builder
func builderConcat(n int, str string) string {
	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString(str)
	}
	return builder.String()
}
// []byte
func byteConcat(n int, str string) string {
	buf := make([]byte, 0)
	for i := 0; i < n; i++ {
		buf = append(buf, str...)
	}
	return string(buf)
}
// bytes.Buffer
func bufferConcat(n int, s string) string {
	buf := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		buf.WriteString(s)
	}
	return buf.String()
}
// +
func plusConcat(n int, str string) string {
	s := ""
	for i := 0; i < n; i++ {
		s += str
	}
	return s
}
func benchmark(b *testing.B, f func(int, string) string) {
	var str = "randomString"
	for i := 0; i < b.N; i++ {
		f(10000, str)
	}
}
func BenchmarkPlusConcat(b *testing.B)    { benchmark(b, plusConcat) }
func BenchmarkBuilderConcat(b *testing.B) { benchmark(b, builderConcat) }
func BenchmarkBufferConcat(b *testing.B)  { benchmark(b, bufferConcat) }
func BenchmarkByteConcat(b *testing.B)    { benchmark(b, byteConcat) }
// 2048 以前按倍数申请，2048 之后，以 640 递增，最后一次递增 24576 到 122880
func TestBuilderConcat(t *testing.T) {
	var str = "random_str"
	var builder strings.Builder
	cap := 0
	for i := 0; i < 10000; i++ {
		if builder.Cap() != cap {
			fmt.Print(builder.Cap(), " ")
			cap = builder.Cap()
		}
		builder.WriteString(str)
	}
}
type Config struct {
	Name    string `json:"server-name"` // CONFIG_SERVER_NAME
	IP      string `json:"server-ip"`   // CONFIG_SERVER_IP
	URL     string `json:"server-url"`  // CONFIG_SERVER_URL
	Timeout string `json:"timeout"`     // CONFIG_TIMEOUT
}
func readConfig() *Config {
	// read from xxx.json，省略
	config := Config{}
	typ := reflect.TypeOf(config)
	value := reflect.Indirect(reflect.ValueOf(&config))
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if v, ok := f.Tag.Lookup("json"); ok {
			key := fmt.Sprintf("CONFIG_%s", strings.ReplaceAll(strings.ToUpper(v), "-", "_"))
			if env, exist := os.LookupEnv(key); exist {
				value.FieldByName(f.Name).Set(reflect.ValueOf(env))
			}
		}
	}
	return &config
}
func TestReadConfig(t *testing.T) {
	os.Setenv("CONFIG_SERVER_NAME", "global_server\n")
	os.Setenv("CONFIG_SERVER_IP", "10.0.0.1\n")
	os.Setenv("CONFIG_SERVER_URL", "geektutu.com\n")
	c := readConfig()
	fmt.Printf("%+v", c)
}
func BenchmarkNew(b *testing.B) {
	var config *Config
	for i := 0; i < b.N; i++ {
		config = new(Config)
	}
	_ = config
}

func BenchmarkReflectNew(b *testing.B) {
	var config *Config
	typ := reflect.TypeOf(Config{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config, _ = reflect.New(typ).Interface().(*Config)
	}
	_ = config
}