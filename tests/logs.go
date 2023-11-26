package tests

func MockLogsData(app *TestApp) error {
	_, err := app.LogsDB().NewQuery(`
		delete from {{_logs}};

		insert into {{_logs}} (
			[[id]],
			[[level]],
			[[message]],
			[[data]],
			[[created]],
			[[updated]]
		)
		values
		(
			"873f2133-9f38-44fb-bf82-c8f53b310d91",
			0,
			"test_message1",
			'{"status":200}',
			"2022-05-01 10:00:00.123Z",
			"2022-05-01 10:00:00.123Z"
		),
		(
			"f2133873-44fb-9f38-bf82-c918f53b310d",
			8,
			"test_message2",
			'{"status":400}',
			"2022-05-02 10:00:00.123Z",
			"2022-05-02 10:00:00.123Z"
		);
	`).Execute()

	return err
}
