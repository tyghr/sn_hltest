wrk.method = "POST"
wrk.body   = "first_name=%25n%25&second_name=%25l%25"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"