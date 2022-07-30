request = function()
    wrk.method = "POST"
    wrk.headers["content-type"] = "application/json"
    username = "test01"
    password = "test01"
    --math.randomseed(tonumber(tostring(os.time()):reverse():sub(1,7)))

    --wrk.headers["Authorization"] = "testToken"
    path = "/api/login?requestID=" .. tostring(math.random(10, 10000000))

    wrk.body = string.format('{"username": "%s", "password": "%s"}', username, password)
    return wrk.format(nil, path)
end
