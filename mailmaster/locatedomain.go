package mailmaster

import (
    . "active_apple/ml/dict"
)

var pubkey = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD30Cx3fgd6wSMpmy71N6L2+70S
1NQwi8JvpQ/ifhcy8M+MFnaU1Zw44FebXSCQGCJf9xIHDSVNi1tvULYPwZPW8NO/
nIYz6JEYZDsyyTazphipJF5eZ01DtQWwFoZgkEf2M6TCOUY56Km3sXPQ1rVhvru/
dMnUNl5PHDuQZbMtNQIDAQAB
-----END PUBLIC KEY-----`

// var pubkey = `
// -----BEGIN PUBLIC KEY-----
// MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCcmGH5Hqbn0aZoGBaqYs0kDAva
// 1D1sccTPoZU+cmVHjOaVtcsM6PYjo5C6MxVHKdwufmozOjxAezXyY5ZxfvvKlAfA
// YFDLPSD6hsOP9NuZb021Gp1vrU8xr4m7p0zKn8Ss4elXVvkl62bJtmAWErPnlVVx
// /CF/OxRMZrzq4JLWMQIDAQAB
// -----END PUBLIC KEY-----`

// var privateKey = `
// -----BEGIN RSA PRIVATE KEY-----
// MIICXAIBAAKBgQCcmGH5Hqbn0aZoGBaqYs0kDAva1D1sccTPoZU+cmVHjOaVtcsM
// 6PYjo5C6MxVHKdwufmozOjxAezXyY5ZxfvvKlAfAYFDLPSD6hsOP9NuZb021Gp1v
// rU8xr4m7p0zKn8Ss4elXVvkl62bJtmAWErPnlVVx/CF/OxRMZrzq4JLWMQIDAQAB
// AoGBAIBJEdIdO0ykYrfaLA9Pu5DxUXDm+J7zoPEcBYDQBIqWMnypHnwoCSTvJWx0
// 1tSixV9NbsEizyNgDLTSwveduLVWzu1oj2IfRipUqXcfXHGG1KGt1Pm0+MqsSl3S
// FE5myql0HEvUNbca8bAHg7d9GGyAZL0JqN/rxo4N4oWDr2hFAkEAu8NZd4Rm4zAu
// aoJ53Gh4suN2IIk181Pq5ehZkh2SyoQ5aDm71RHMDcAqZutuXJrlh/JDo84db/zp
// mkZctiEo/wJBANWBUJKiGrVE9nUk37cJlFlqtzM6QQz6sRR6B9enq41SuL58z+Ss
// pyj1+Y/bMf6t/0cvua31USFP/q2769YNUM8CQBT4xM1spHFLuGN9H09W++Q/M7p+
// mOAMx3fWc+q2Euc7zY2upSQvULNYe2Pzd+gwBOMiVBu/sdoITa9FnKVbHtECQAr4
// 4lMm0YiPSrsqcfTOITmXKmMPk1g/aepLeyuyCjbxEV14vJZb6RtJyNGDykX0WzIl
// Wb1+5fR4T/ZNugj+FjECQH9XAkTYgHEtpCrZGcoUvIJ2hagOSBOweMNwIYBZZBzT
// EuI2FO2dOYhcjEBbsldx7yx4/NrQp0ReCYtTFzSLNQg=
// -----END RSA PRIVATE KEY-----`

var locateDomain = Dict{
    "locate_domain"     : "lbs.client.163.com",
    "locate_port"       : 8080,
    "locate_addr"       : "192.168.130.100",
    "pubkey"            : pubkey,
    "pubkey_version"    : 1,
}
