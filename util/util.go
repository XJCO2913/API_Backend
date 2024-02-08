package util

import "reflect"

func IsEmpty(value interface{}) bool {
    if value == nil {
        return true
    }

    v := reflect.ValueOf(value)
    switch v.Kind() {
    case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
        return v.Len() == 0
    case reflect.Bool:
        return !v.Bool()
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return v.Int() == 0
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
        return v.Uint() == 0
    case reflect.Float32, reflect.Float64:
        return v.Float() == 0
    case reflect.Interface, reflect.Ptr:
        return v.IsNil()
    }

    return false
}
