local key = KEYS[1]
local delta = tonumber(ARGV[1])
local exp = tonumber(ARGV[2])
local exists = redis.call("EXISTS", key)
if exists == 1 then
    -- 说明key存在
    redis.call("INCRBY", key, delta)
    -- 设置过期时间
    redis.call("EXPIRE", key, exp)
else
    -- 说明key不存在
    return
end