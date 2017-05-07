package json

import (
    . "ml/trace"
    . "ml/dict"

    "os"
    "io/ioutil"
)

func MustMarshal(v interface{}) []byte {
    data, err := Marshal(v)
    RaiseIf(err)
    return data
}

func MustMarshalIndent(v interface{}, prefix, indent string) []byte {
    data, err := MarshalIndent(v, prefix, indent)
    RaiseIf(err)
    return data
}

func MustUnmarshal(data []byte, v interface{}) {
    err := Unmarshal(data, v)
    if err != nil {
        Raise(NewJSONDecodeError(err.Error()))
    }
}

func LoadFile(file string, v interface{}) {
    f, err := os.Open(file)
    if err != nil {
        Raise(NewFileNotFoundError(err.Error()))
    }

    defer f.Close()

    data, err := ioutil.ReadAll(f)
    if err != nil {
        Raise(NewBaseException(err.Error()))
    }

    LoadData(data, v)
}

func LoadData(data []byte, v interface{}) {
    err := Unmarshal(data, v)
    if err != nil {
        Raise(NewJSONDecodeError(err.Error()))
    }
}

func LoadString(text string, v interface{}) {
    LoadData([]byte(text), v)
}

func LoadFileDict(file string) (v JsonDict) {
    LoadFile(file, &v)
    return
}

func LoadDataDict(data []byte) (v JsonDict) {
    LoadData(data, &v)
    return
}

func LoadStringDict(text string) (v JsonDict) {
    LoadString(text, &v)
    return
}

func LoadFileArray(file string) (v JsonArray) {
    LoadFile(file, &v)
    return
}

func LoadDataArray(data []byte) (v JsonArray) {
    LoadData(data, &v)
    return
}

func LoadStringArray(text string) (v JsonArray) {
    LoadString(text, &v)
    return
}
