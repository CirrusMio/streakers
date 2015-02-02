-- wrk script that loops over the user.list file and makes POST requests for each user.

file = assert(io.open("user.list"))

request = function()
  wrk.method = "POST"
  wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
  wrk.body = "name=" .. file:read()
  return wrk.format()
end

done = function(summary, latency, requests)
  file:close()
end
