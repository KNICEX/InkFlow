-- 验证码的key
-- phone_code:login:173xxx
local key = KEYS[1]
-- 构造验证次数对应的key
-- phone_code:login:173xxx:cnt
local cntKey = key..":cnt"
local val = ARGV[1]
-- 验证码有效期
local ex = tonumber(ARGV[2])
-- 重新发送间隔
local resendInterval = tonumber(ARGV[3])
-- 最大验证次数
local maxRetry = tonumber(ARGV[4])
-- 验证码有效时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key存在,但是没有过期时间,异常状态
    return -2
elseif ttl == -2 or ttl < (ex - resendInterval) then
    -- 没有key 或者 key存在时间已经大于1分钟(10-1 minutes=540s)
    -- setex key timeout val
    redis.call("setex", key, ex, val)
    redis.call("setex", cntKey, ex, maxRetry)
    return 0
else
    -- 发送太频繁
    return -1
end
