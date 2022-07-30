request = function()
    wrk.method = "GET"
    wrk.headers["content-type"] = "application/json"
    --username = "test01"
    username = "test" .. tostring(math.random(10, 10000000))

    wrk.headers["Authorization"] = "testToken"
    path = "/api/profile?username=" .. username .. "&requestID=" .. tostring(math.random(10, 10000000))
    return wrk.format(nil, path)
end
