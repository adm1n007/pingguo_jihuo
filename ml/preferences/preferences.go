package preferences

import (
    "encoding/json"
    "os"
    "io/ioutil"
)

func LoadString(str string, v interface{}) error {
    return json.Unmarshal([]byte(str), v)
}

func LoadBytes(buf []byte, v interface{}) error {
    return json.Unmarshal(buf, v)
}

func LoadFile(fileName string, v interface{}) error {
    file, err := os.Open(fileName)
    if err != nil {
        return err
    }

    defer file.Close()

    buf, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }

    return json.Unmarshal(buf, v)
}
