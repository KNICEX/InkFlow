local key = KEYS[1]
local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key .. ":cnt"
local cnt = tonumber(redis.call("get", cntKey))
if cnt == nil or cnt <= 0 then
    --   验证次数用完/验证码过期/验证码被使用过/验证码不存在
    return -1
elseif expectedCode == code  then
    --   验证码正确
    redis.call("del", key)
    redis.call("del", cntKey)
    return 0
else
    --   验证码错误
    redis.call("decr", cntKey)
    return -2
end