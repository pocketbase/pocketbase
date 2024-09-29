package tests

import (
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func StubOTPRecords(app core.App) error {
	superuser2, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test2@example.com")
	if err != nil {
		return err
	}
	superuser2.SetRaw("stubId", "superuser2")

	superuser3, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test3@example.com")
	if err != nil {
		return err
	}
	superuser3.SetRaw("stubId", "superuser3")

	user1, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		return err
	}
	user1.SetRaw("stubId", "user1")

	now := types.NowDateTime()
	old := types.NowDateTime().Add(-1 * time.Hour)

	stubs := map[*core.Record][]types.DateTime{
		superuser2: {now, now.Add(-1 * time.Millisecond), old, now.Add(-2 * time.Millisecond), old.Add(-1 * time.Millisecond)},
		superuser3: {now.Add(-3 * time.Millisecond), now.Add(-2 * time.Minute)},
		user1:      {old},
	}
	for record, idDates := range stubs {
		for i, date := range idDates {
			otp := core.NewOTP(app)
			otp.Id = record.GetString("stubId") + "_" + strconv.Itoa(i)
			otp.SetRecordRef(record.Id)
			otp.SetCollectionRef(record.Collection().Id)
			otp.SetPassword("test123")
			otp.SetRaw("created", date)
			if err := app.SaveNoValidate(otp); err != nil {
				return err
			}
		}
	}

	return nil
}

func StubMFARecords(app core.App) error {
	superuser2, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test2@example.com")
	if err != nil {
		return err
	}
	superuser2.SetRaw("stubId", "superuser2")

	superuser3, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test3@example.com")
	if err != nil {
		return err
	}
	superuser3.SetRaw("stubId", "superuser3")

	user1, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		return err
	}
	user1.SetRaw("stubId", "user1")

	now := types.NowDateTime()
	old := types.NowDateTime().Add(-1 * time.Hour)

	type mfaData struct {
		method string
		date   types.DateTime
	}

	stubs := map[*core.Record][]mfaData{
		superuser2: {
			{core.MFAMethodOTP, now},
			{core.MFAMethodOTP, old},
			{core.MFAMethodPassword, now.Add(-2 * time.Minute)},
			{core.MFAMethodPassword, now.Add(-1 * time.Millisecond)},
			{core.MFAMethodOAuth2, old.Add(-1 * time.Millisecond)},
		},
		superuser3: {
			{core.MFAMethodOAuth2, now.Add(-3 * time.Millisecond)},
			{core.MFAMethodPassword, now.Add(-3 * time.Minute)},
		},
		user1: {
			{core.MFAMethodOAuth2, old},
		},
	}
	for record, idDates := range stubs {
		for i, data := range idDates {
			otp := core.NewMFA(app)
			otp.Id = record.GetString("stubId") + "_" + strconv.Itoa(i)
			otp.SetRecordRef(record.Id)
			otp.SetCollectionRef(record.Collection().Id)
			otp.SetMethod(data.method)
			otp.SetRaw("created", data.date)
			if err := app.SaveNoValidate(otp); err != nil {
				return err
			}
		}
	}

	return nil
}

func StubLogsData(app *TestApp) error {
	_, err := app.AuxDB().NewQuery(`
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
