namespace go rpc.xjco2913.auth

include "base.thrift"

struct LoginReq {
    1: required string username
    2: required string password
}

struct LoginResp {
    1: string token
    2: string username
    3: i32 gender
    4: string birthday
    5: string region

    255: base.BaseResp baseResp
}

struct RefreshTokenReq {
    1: required string oldToken
}

struct RefreshTokenResp {
    1: string newToken
}

service AuthService {
    LoginResp Login(1: LoginReq req)
    RefreshTokenResp RefreshToken(1: RefreshTokenReq req)
}