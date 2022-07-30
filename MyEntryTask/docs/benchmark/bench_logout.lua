request = function()
    wrk.method = "POST"
    wrk.headers["content-type"] = "application/json"
    --username = "test01"
    username = "test" .. tostring(math.random(10, 10000000))
    wrk.headers["Authorization"] = "testToken"
    path = "/api/signout?requestID=" .. tostring(math.random(10, 10000000))

    wrk.body = string.format('{"username": "%s", "password": "%s"}', username, password)
    return wrk.format(nil, path)
end