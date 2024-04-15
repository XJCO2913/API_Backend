namespace go rpc.xjco2913.base

struct BaseResp {
    1: i16 code
    2: string msg
    3: map<string, string> data
}