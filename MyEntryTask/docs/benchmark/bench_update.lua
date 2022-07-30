request = function()
    wrk.method = "POST"
    wrk.headers["content-type"] = "application/x-www-form-urlencoded"
    username = "test" .. math.random(10, 10000000)
    --username = "test01"
    nickname = "测试011"
    wrk.headers["Authorization"] = "testToken"
    body = "username=" .. username .. "&nickname=" .. nickname
    path = "/api/updateProfile?requestID=" .. tostring(math.random(10, 10000000))

    return wrk.format(nil, path, nil, body)
end