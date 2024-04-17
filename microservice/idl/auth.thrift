namespace go rpc.xjco2913.auth

include "base.thrift"

struct LoginReq {
    1: required string username
    2: required string password
}

struct LoginResp {
    1: string token
    2: string username
    3: string gender
    4: string birthday
    5: string region

    255: base.BaseResp baseResp
}

service LoginService {
    LoginResp Login(1: LoginReq req)
}