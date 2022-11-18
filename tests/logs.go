package tests

func MockRequestLogsData(app *TestApp) error {
	_, err := app.LogsDB().NewQuery(`
		delete from {{_requests}};

		insert into {{_requests}} (
			[[id]],
			[[url]],
			[[method]],
			[[status]],
			[[auth]],
			[[userIp]],
			[[remoteIp]],
			[[referer]],
			[[userAgent]],
			[[meta]],
			[[created]],
			[[updated]]
		)
		values
		(
			"873f2133-9f38-44fb-bf82-c8f53b310d91",
			"/test1",
			"get",
			200,
			"guest",
			"127.0.0.1",
			"127.0.0.1",
			"",
			"",
			"{}",
			"2022-05-01 10:00:00.123Z",
			"2022-05-01 10:00:00.123Z"
		),
		(
			"f2133873-44fb-9f38-bf82-c918f53b310d",
			"/test2",
			"post",
			400,
			"admin",
			"127.0.0.1",
			"127.0.0.1",
			"",
			"",
			'{"errorDetails":"error_details..."}',
			"2022-05-02 10:00:00.123Z",
			"2022-05-02 10:00:00.123Z"
		);
	`).Execute()

	return err
}
