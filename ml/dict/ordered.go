package dict

import (
    . "fmt"
    . "active_apple/ml/array"
    . "active_apple/ml/trace"

    "encoding/json"

    "active_apple/ml/io2/filestream"
)

type orderedDict struct {
    dict Dict
    keys Array
}

type OrderedDict struct {
    *orderedDict
}

func NewOrderedDict(items ...DictItem) OrderedDict {
    d := OrderedDict{
        orderedDict : &orderedDict{
            dict: Dict{},
            keys: Array{},
        },
    }

    for _, item := range items {
        d.Set(item.Key, item.Value)
    }

    return d
}

func (self *orderedDict) toString(depth int) string {
    space := ""

    for i := depth; i > 0; i-- {
        space += "  "
    }

    s := "{\n"
    for _, k := range self.Keys() {
        var key, value string

        v := self.Get(k)

        switch t := k.(type) {
            case string:
                key = Sprintf("%q", t)

            default:
                key = Sprintf("%v", t)
        }

        switch obj := v.(type) {
            case *orderedDict:
                value = obj.toString(depth + 1)

            case orderedDict:
                value = obj.toString(depth + 1)

            case Dict:
                value = obj.toString(depth + 1)

            case string:
                value = Sprintf("%q", obj)

            default:
                value = Sprintf("%+v", obj)
        }

        s += Sprintf("%v  %+v: %+v,\n", space, key, value)
    }

    s += space + "}"

    return s
}

func (self *orderedDict) String() string {
    return self.toString(0)
    // return self.dict.String()
}

func (self *orderedDict) Dict() Dict {
    return self.dict
}

func (self *orderedDict) Length() int {
    return len(self.keys)
}

func (self *orderedDict) Clear() {
    for _, item := range self.Items() {
        self.Remove(item.Key)
    }
}

func (self *orderedDict) Get(key interface{}) interface{} {
    return self.dict[key]
}

func (self *orderedDict) Set(key, value interface{}) {
    if self.keys.Contains(key) == false {
        self.keys.Append(key)
    }

    self.dict[key] = value
}

func (self *orderedDict) Remove(key interface{}) {
    delete(self.dict, key)
    self.keys.Remove(key)
}

func (self *orderedDict) Keys() Array {
    return self.keys
}

func (self *orderedDict) Values() Array {
    values := Array{}

    for _, key := range self.keys {
        values.Append(self.Get(key))
    }

    return values
}

func (self *orderedDict) Items() []DictItem {
    items := make([]DictItem, self.keys.Length())

    for index, key := range self.keys {
        items[index] = DictItem{Key: key, Value: self.Get(key)}
    }

    return items
}

func (self *orderedDict) MergeFrom(other OrderedDict) *orderedDict {
    for _, k := range other.Keys() {
        self.Set(k, other.Get(k))
    }

    return self
}

func (self *orderedDict) MarshalJSON() ([]byte, error) {
    buffer := filestream.CreateMemory()

    buffer.Write("{")

    for i, item := range self.Items() {
        if i != 0 {
            buffer.Write(byte(','))
        }

        buffer.Write(Sprintf("%q:", item.Key))
        data, err := json.Marshal(item.Value)
        RaiseIf(err)
        buffer.Write(data)
    }

    buffer.Write("}")

    return buffer.ReadAll(), nil
}

func (self *orderedDict) JsonIndent(indent string) []byte {
    data, err := json.MarshalIndent(self, "", indent)
    RaiseIf(err)

    return data
}

func (self *orderedDict) Json(indent ...string) []byte {
    var data []byte
    var err error

    switch len(indent) {
        case 0:
            data, err = json.Marshal(self)

        default:
            data, err = json.MarshalIndent(self, "", indent[0])
    }

    RaiseIf(err)

    return data
}
