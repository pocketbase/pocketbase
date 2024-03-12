package jsvm

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dop251/goja"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

func testBindsCount(vm *goja.Runtime, namespace string, count int, t *testing.T) {
	v, err := vm.RunString(`Object.keys(` + namespace + `).length`)
	if err != nil {
		t.Fatal(err)
	}

	total, _ := v.Export().(int64)

	if int(total) != count {
		t.Fatalf("Expected %d %s binds, got %d", count, namespace, total)
	}
}

// note: this test is useful as a reminder to update the tests in case
// a new base binding is added.
func TestBaseBindsCount(t *testing.T) {
	vm := goja.New()
	baseBinds(vm)

	testBindsCount(vm, "this", 17, t)
}

func TestBaseBindsSleep(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	baseBinds(vm)
	vm.Set("reader", strings.NewReader("test"))

	start := time.Now()
	_, err := vm.RunString(`
		sleep(100);
	`)
	if err != nil {
		t.Fatal(err)
	}

	lasted := time.Since(start).Milliseconds()
	if lasted < 100 || lasted > 150 {
		t.Fatalf("Expected to sleep for ~100ms, got %d", lasted)
	}
}

func TestBaseBindsReaderToString(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	baseBinds(vm)
	vm.Set("reader", strings.NewReader("test"))

	_, err := vm.RunString(`
		let result = readerToString(reader)

		if (result != "test") {
			throw new Error('Expected "test", got ' + result);
		}
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBaseBindsCookie(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	baseBinds(vm)

	_, err := vm.RunString(`
		const cookie = new Cookie({
			name:     "example_name",
			value:    "example_value",
			path:     "/example_path",
			domain:   "example.com",
			maxAge:   10,
			secure:   true,
			httpOnly: true,
			sameSite: 3,
		});

		const result = cookie.string();

		const expected = "example_name=example_value; Path=/example_path; Domain=example.com; Max-Age=10; HttpOnly; Secure; SameSite=Strict";

		if (expected != result) {
			throw new("Expected \n" + expected + "\ngot\n" + result);
		}
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBaseBindsSubscriptionMessage(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	baseBinds(vm)
	vm.Set("bytesToString", func(b []byte) string {
		return string(b)
	})

	_, err := vm.RunString(`
		const payload = {
			name: "test",
			data: '{"test":123}'
		}

		const result = new SubscriptionMessage(payload);

		if (result.name != payload.name) {
			throw new("Expected name " + payload.name + ", got " + result.name);
		}

		if (bytesToString(result.data) != payload.data) {
			throw new("Expected data '" + payload.data + "', got '" + bytesToString(result.data) + "'");
		}
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBaseBindsRecord(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	vm := goja.New()
	baseBinds(vm)
	vm.Set("collection", collection)

	// without record data
	// ---
	v1, err := vm.RunString(`new Record(collection)`)
	if err != nil {
		t.Fatal(err)
	}

	m1, ok := v1.Export().(*models.Record)
	if !ok {
		t.Fatalf("Expected m1 to be models.Record, got \n%v", m1)
	}

	// with record data
	// ---
	v2, err := vm.RunString(`new Record(collection, { email: "test@example.com" })`)
	if err != nil {
		t.Fatal(err)
	}

	m2, ok := v2.Export().(*models.Record)
	if !ok {
		t.Fatalf("Expected m2 to be models.Record, got \n%v", m2)
	}

	if m2.Collection().Name != "users" {
		t.Fatalf("Expected record with collection %q, got \n%v", "users", m2.Collection())
	}

	if m2.Email() != "test@example.com" {
		t.Fatalf("Expected record with email field set to %q, got \n%v", "test@example.com", m2)
	}
}

func TestBaseBindsCollection(t *testing.T) {
	vm := goja.New()
	baseBinds(vm)

	v, err := vm.RunString(`new Collection({ name: "test", createRule: "@request.auth.id != ''", schema: [{name: "title", "type": "text"}] })`)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := v.Export().(*models.Collection)
	if !ok {
		t.Fatalf("Expected models.Collection, got %v", m)
	}

	if m.Name != "test" {
		t.Fatalf("Expected collection with name %q, got %q", "test", m.Name)
	}

	expectedRule := "@request.auth.id != ''"
	if m.CreateRule == nil || *m.CreateRule != expectedRule {
		t.Fatalf("Expected create rule %q, got %v", "@request.auth.id != ''", m.CreateRule)
	}

	if f := m.Schema.GetFieldByName("title"); f == nil {
		t.Fatalf("Expected schema to be set, got %v", m.Schema)
	}
}

func TestBaseVMAdminBind(t *testing.T) {
	vm := goja.New()
	baseBinds(vm)

	v, err := vm.RunString(`new Admin({ email: "test@example.com" })`)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := v.Export().(*models.Admin)
	if !ok {
		t.Fatalf("Expected models.Admin, got %v", m)
	}
}

func TestBaseBindsSchema(t *testing.T) {
	vm := goja.New()
	baseBinds(vm)

	v, err := vm.RunString(`new Schema([{name: "title", "type": "text"}])`)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := v.Export().(*schema.Schema)
	if !ok {
		t.Fatalf("Expected schema.Schema, got %v", m)
	}

	if f := m.GetFieldByName("title"); f == nil {
		t.Fatalf("Expected schema fields to be loaded, got %v", m.Fields())
	}
}

func TestBaseBindsSchemaField(t *testing.T) {
	vm := goja.New()
	baseBinds(vm)

	v, err := vm.RunString(`new SchemaField({name: "title", "type": "text"})`)
	if err != nil {
		t.Fatal(err)
	}

	f, ok := v.Export().(*schema.SchemaField)
	if !ok {
		t.Fatalf("Expected schema.SchemaField, got %v", f)
	}

	if f.Name != "title" {
		t.Fatalf("Expected field %q, got %v", "title", f)
	}
}

func TestBaseBindsMailerMessage(t *testing.T) {
	vm := goja.New()
	baseBinds(vm)

	v, err := vm.RunString(`new MailerMessage({
		from: {name: "test_from", address: "test_from@example.com"},
		to: [
			{name: "test_to1", address: "test_to1@example.com"},
			{name: "test_to2", address: "test_to2@example.com"},
		],
		bcc: [
			{name: "test_bcc1", address: "test_bcc1@example.com"},
			{name: "test_bcc2", address: "test_bcc2@example.com"},
		],
		cc: [
			{name: "test_cc1", address: "test_cc1@example.com"},
			{name: "test_cc2", address: "test_cc2@example.com"},
		],
		subject: "test_subject",
		html: "test_html",
		text: "test_text",
		headers: {
			header1: "a",
			header2: "b",
		}
	})`)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := v.Export().(*mailer.Message)
	if !ok {
		t.Fatalf("Expected mailer.Message, got %v", m)
	}

	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"from":{"Name":"test_from","Address":"test_from@example.com"},"to":[{"Name":"test_to1","Address":"test_to1@example.com"},{"Name":"test_to2","Address":"test_to2@example.com"}],"bcc":[{"Name":"test_bcc1","Address":"test_bcc1@example.com"},{"Name":"test_bcc2","Address":"test_bcc2@example.com"}],"cc":[{"Name":"test_cc1","Address":"test_cc1@example.com"},{"Name":"test_cc2","Address":"test_cc2@example.com"}],"subject":"test_subject","html":"test_html","text":"test_text","headers":{"header1":"a","header2":"b"},"attachments":null}`

	if string(raw) != expected {
		t.Fatalf("Expected \n%s, \ngot \n%s", expected, raw)
	}
}

func TestBaseBindsCommand(t *testing.T) {
	vm := goja.New()
	baseBinds(vm)

	_, err := vm.RunString(`
		let runCalls = 0;

		let cmd = new Command({
			use: "test",
			run: (c, args) => {
				runCalls++;
			}
		});

		cmd.run(null, []);

		if (cmd.use != "test") {
			throw new Error('Expected cmd.use "test", got: ' + cmd.use);
		}

		if (runCalls != 1) {
			throw new Error('Expected runCalls 1, got: ' + runCalls);
		}
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBaseBindsRequestInfo(t *testing.T) {
	vm := goja.New()
	baseBinds(vm)

	_, err := vm.RunString(`
		let info = new RequestInfo({
			admin: new Admin({id: "test1"}),
			data: {"name": "test2"}
		});

		if (info.admin?.id != "test1") {
			throw new Error('Expected info.admin.id to be test1, got: ' + info.admin?.id);
		}

		if (info.data?.name != "test2") {
			throw new Error('Expected info.data.name to be test2, got: ' + info.data?.name);
		}
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBaseBindsDateTime(t *testing.T) {
	vm := goja.New()
	baseBinds(vm)

	_, err := vm.RunString(`
		const v0 = new DateTime();
		if (v0.isZero()) {
			throw new Error('Expected to fallback to now, got zero value');
		}

		const v1 = new DateTime('2023-01-01 00:00:00.000Z');
		const expected = "2023-01-01 00:00:00.000Z"
		if (v1.string() != expected) {
			throw new Error('Expected ' + expected + ', got ', v1.string());
		}
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBaseBindsValidationError(t *testing.T) {
	vm := goja.New()
	baseBinds(vm)

	scenarios := []struct {
		js            string
		expectCode    string
		expectMessage string
	}{
		{
			`new ValidationError()`,
			"",
			"",
		},
		{
			`new ValidationError("test_code")`,
			"test_code",
			"",
		},
		{
			`new ValidationError("test_code", "test_message")`,
			"test_code",
			"test_message",
		},
	}

	for _, s := range scenarios {
		v, err := vm.RunString(s.js)
		if err != nil {
			t.Fatal(err)
		}

		m, ok := v.Export().(validation.Error)
		if !ok {
			t.Fatalf("[%s] Expected validation.Error, got %v", s.js, m)
		}

		if m.Code() != s.expectCode {
			t.Fatalf("[%s] Expected code %q, got %q", s.js, s.expectCode, m.Code())
		}

		if m.Message() != s.expectMessage {
			t.Fatalf("[%s] Expected message %q, got %q", s.js, s.expectMessage, m.Message())
		}
	}
}

func TestBaseBindsDao(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	baseBinds(vm)
	vm.Set("db", app.Dao().ConcurrentDB())
	vm.Set("db2", app.Dao().NonconcurrentDB())

	scenarios := []struct {
		js              string
		concurrentDB    dbx.Builder
		nonconcurrentDB dbx.Builder
	}{
		{
			js:              "new Dao(db)",
			concurrentDB:    app.Dao().ConcurrentDB(),
			nonconcurrentDB: app.Dao().ConcurrentDB(),
		},
		{
			js:              "new Dao(db, db2)",
			concurrentDB:    app.Dao().ConcurrentDB(),
			nonconcurrentDB: app.Dao().NonconcurrentDB(),
		},
	}

	for _, s := range scenarios {
		v, err := vm.RunString(s.js)
		if err != nil {
			t.Fatalf("[%s] Failed to execute js script, got %v", s.js, err)
		}

		d, ok := v.Export().(*daos.Dao)
		if !ok {
			t.Fatalf("[%s] Expected daos.Dao, got %v", s.js, d)
		}

		if d.ConcurrentDB() != s.concurrentDB {
			t.Fatalf("[%s] The ConcurrentDB instances doesn't match", s.js)
		}

		if d.NonconcurrentDB() != s.nonconcurrentDB {
			t.Fatalf("[%s] The NonconcurrentDB instances doesn't match", s.js)
		}
	}
}

func TestDbxBinds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	vm.Set("db", app.Dao().DB())
	baseBinds(vm)
	dbxBinds(vm)

	testBindsCount(vm, "$dbx", 15, t)

	sceneraios := []struct {
		js       string
		expected string
	}{
		{
			`$dbx.exp("a = 1").build(db, {})`,
			"a = 1",
		},
		{
			`$dbx.hashExp({
				"a": 1,
				b: null,
				c: [1, 2, 3],
			}).build(db, {})`,
			"`a`={:p0} AND `b` IS NULL AND `c` IN ({:p1}, {:p2}, {:p3})",
		},
		{
			`$dbx.not($dbx.exp("a = 1")).build(db, {})`,
			"NOT (a = 1)",
		},
		{
			`$dbx.and($dbx.exp("a = 1"), $dbx.exp("b = 2")).build(db, {})`,
			"(a = 1) AND (b = 2)",
		},
		{
			`$dbx.or($dbx.exp("a = 1"), $dbx.exp("b = 2")).build(db, {})`,
			"(a = 1) OR (b = 2)",
		},
		{
			`$dbx.in("a", 1, 2, 3).build(db, {})`,
			"`a` IN ({:p0}, {:p1}, {:p2})",
		},
		{
			`$dbx.notIn("a", 1, 2, 3).build(db, {})`,
			"`a` NOT IN ({:p0}, {:p1}, {:p2})",
		},
		{
			`$dbx.like("a", "test1", "test2").match(true, false).build(db, {})`,
			"`a` LIKE {:p0} AND `a` LIKE {:p1}",
		},
		{
			`$dbx.orLike("a", "test1", "test2").match(false, true).build(db, {})`,
			"`a` LIKE {:p0} OR `a` LIKE {:p1}",
		},
		{
			`$dbx.notLike("a", "test1", "test2").match(true, false).build(db, {})`,
			"`a` NOT LIKE {:p0} AND `a` NOT LIKE {:p1}",
		},
		{
			`$dbx.orNotLike("a", "test1", "test2").match(false, false).build(db, {})`,
			"`a` NOT LIKE {:p0} OR `a` NOT LIKE {:p1}",
		},
		{
			`$dbx.exists($dbx.exp("a = 1")).build(db, {})`,
			"EXISTS (a = 1)",
		},
		{
			`$dbx.notExists($dbx.exp("a = 1")).build(db, {})`,
			"NOT EXISTS (a = 1)",
		},
		{
			`$dbx.between("a", 1, 2).build(db, {})`,
			"`a` BETWEEN {:p0} AND {:p1}",
		},
		{
			`$dbx.notBetween("a", 1, 2).build(db, {})`,
			"`a` NOT BETWEEN {:p0} AND {:p1}",
		},
	}

	for _, s := range sceneraios {
		result, err := vm.RunString(s.js)
		if err != nil {
			t.Fatalf("[%s] Failed to execute js script, got %v", s.js, err)
		}

		v, _ := result.Export().(string)

		if v != s.expected {
			t.Fatalf("[%s] Expected \n%s, \ngot \n%s", s.js, s.expected, v)
		}
	}
}

func TestMailsBindsCount(t *testing.T) {
	vm := goja.New()
	mailsBinds(vm)

	testBindsCount(vm, "$mails", 4, t)
}

func TestMailsBinds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	admin, err := app.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	record, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	vm := goja.New()
	baseBinds(vm)
	mailsBinds(vm)
	vm.Set("$app", app)
	vm.Set("admin", admin)
	vm.Set("record", record)

	_, vmErr := vm.RunString(`
		$mails.sendAdminPasswordReset($app, admin);
		if (!$app.testMailer.lastMessage.html.includes("/_/#/confirm-password-reset/")) {
			throw new Error("Expected admin password reset email")
		}

		$mails.sendRecordPasswordReset($app, record);
		if (!$app.testMailer.lastMessage.html.includes("/_/#/auth/confirm-password-reset/")) {
			throw new Error("Expected record password reset email")
		}

		$mails.sendRecordVerification($app, record);
		if (!$app.testMailer.lastMessage.html.includes("/_/#/auth/confirm-verification/")) {
			throw new Error("Expected record verification email")
		}

		$mails.sendRecordChangeEmail($app, record, "new@example.com");
		if (!$app.testMailer.lastMessage.html.includes("/_/#/auth/confirm-email-change/")) {
			throw new Error("Expected record email change email")
		}
	`)
	if vmErr != nil {
		t.Fatal(vmErr)
	}
}

func TestTokensBindsCount(t *testing.T) {
	vm := goja.New()
	tokensBinds(vm)

	testBindsCount(vm, "$tokens", 8, t)
}

func TestTokensBinds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	admin, err := app.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	record, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	vm := goja.New()
	baseBinds(vm)
	tokensBinds(vm)
	vm.Set("$app", app)
	vm.Set("admin", admin)
	vm.Set("record", record)

	sceneraios := []struct {
		js  string
		key string
	}{
		{
			`$tokens.adminAuthToken($app, admin)`,
			admin.TokenKey + app.Settings().AdminAuthToken.Secret,
		},
		{
			`$tokens.adminResetPasswordToken($app, admin)`,
			admin.TokenKey + app.Settings().AdminPasswordResetToken.Secret,
		},
		{
			`$tokens.adminFileToken($app, admin)`,
			admin.TokenKey + app.Settings().AdminFileToken.Secret,
		},
		{
			`$tokens.recordAuthToken($app, record)`,
			record.TokenKey() + app.Settings().RecordAuthToken.Secret,
		},
		{
			`$tokens.recordVerifyToken($app, record)`,
			record.TokenKey() + app.Settings().RecordVerificationToken.Secret,
		},
		{
			`$tokens.recordResetPasswordToken($app, record)`,
			record.TokenKey() + app.Settings().RecordPasswordResetToken.Secret,
		},
		{
			`$tokens.recordChangeEmailToken($app, record)`,
			record.TokenKey() + app.Settings().RecordEmailChangeToken.Secret,
		},
		{
			`$tokens.recordFileToken($app, record)`,
			record.TokenKey() + app.Settings().RecordFileToken.Secret,
		},
	}

	for _, s := range sceneraios {
		result, err := vm.RunString(s.js)
		if err != nil {
			t.Fatalf("[%s] Failed to execute js script, got %v", s.js, err)
		}

		v, _ := result.Export().(string)

		if _, err := security.ParseJWT(v, s.key); err != nil {
			t.Fatalf("[%s] Failed to parse JWT %v, got %v", s.js, v, err)
		}
	}
}

func TestSecurityBindsCount(t *testing.T) {
	vm := goja.New()
	securityBinds(vm)

	testBindsCount(vm, "$security", 15, t)
}

func TestSecurityCryptoBinds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	baseBinds(vm)
	securityBinds(vm)

	sceneraios := []struct {
		js       string
		expected string
	}{
		{`$security.md5("123")`, "202cb962ac59075b964b07152d234b70"},
		{`$security.sha256("123")`, "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"},
		{`$security.sha512("123")`, "3c9909afec25354d551dae21590bb26e38d53f2173b8d3dc3eee4c047e7ab1c1eb8b85103e3be7ba613b31bb5c9c36214dc9f14a42fd7a2fdb84856bca5c44c2"},
		{`$security.hs256("hello", "test")`, "f151ea24bda91a18e89b8bb5793ef324b2a02133cce15a28a719acbd2e58a986"},
		{`$security.hs512("hello", "test")`, "44f280e11103e295c26cd61dd1cdd8178b531b860466867c13b1c37a26b6389f8af110efbe0bb0717b9d9c87f6fe1c97b3b1690936578890e5669abf279fe7fd"},
		{`$security.equal("abc", "abc")`, "true"},
		{`$security.equal("abc", "abcd")`, "false"},
	}

	for _, s := range sceneraios {
		t.Run(s.js, func(t *testing.T) {
			result, err := vm.RunString(s.js)
			if err != nil {
				t.Fatalf("Failed to execute js script, got %v", err)
			}

			v := cast.ToString(result.Export())

			if v != s.expected {
				t.Fatalf("Expected %v \ngot \n%v", s.expected, v)
			}
		})
	}
}

func TestSecurityRandomStringBinds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	baseBinds(vm)
	securityBinds(vm)

	sceneraios := []struct {
		js     string
		length int
	}{
		{`$security.randomString(6)`, 6},
		{`$security.randomStringWithAlphabet(7, "abc")`, 7},
		{`$security.pseudorandomString(8)`, 8},
		{`$security.pseudorandomStringWithAlphabet(9, "abc")`, 9},
	}

	for _, s := range sceneraios {
		t.Run(s.js, func(t *testing.T) {
			result, err := vm.RunString(s.js)
			if err != nil {
				t.Fatalf("Failed to execute js script, got %v", err)
			}

			v, _ := result.Export().(string)

			if len(v) != s.length {
				t.Fatalf("Expected %d length string, \ngot \n%v", s.length, v)
			}
		})
	}
}

func TestSecurityJWTBinds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	sceneraios := []struct {
		name string
		js   string
	}{
		{
			"$security.parseUnverifiedJWT",
			`
				const result = $security.parseUnverifiedJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIn0.aXzC7q7z1lX_hxk5P0R368xEU7H1xRwnBQQcLAmG0EY")
				if (result.name != "John Doe") {
					throw new Error("Expected result.name 'John Doe', got " + result.name)
				}
				if (result.sub != "1234567890") {
					throw new Error("Expected result.sub '1234567890', got " + result.sub)
				}
			`,
		},
		{
			"$security.parseJWT",
			`
				const result = $security.parseJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIn0.aXzC7q7z1lX_hxk5P0R368xEU7H1xRwnBQQcLAmG0EY", "test")
				if (result.name != "John Doe") {
					throw new Error("Expected result.name 'John Doe', got " + result.name)
				}
				if (result.sub != "1234567890") {
					throw new Error("Expected result.sub '1234567890', got " + result.sub)
				}
			`,
		},
		{
			"$security.createJWT",
			`
				// overwrite the exp claim for static token
				const result = $security.createJWT({"exp": 123}, "test", 0)

				const expected = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEyM30.7gbv7w672gApdBRASI6OniCtKwkKjhieSxsr6vxSrtw";
				if (result != expected) {
					throw new Error("Expected token \n" + expected + ", got \n" + result)
				}
			`,
		},
	}

	for _, s := range sceneraios {
		t.Run(s.name, func(t *testing.T) {
			vm := goja.New()
			baseBinds(vm)
			securityBinds(vm)

			_, err := vm.RunString(s.js)
			if err != nil {
				t.Fatalf("Failed to execute js script, got %v", err)
			}
		})
	}
}

func TestSecurityEncryptAndDecryptBinds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	baseBinds(vm)
	securityBinds(vm)

	_, err := vm.RunString(`
		const key = "abcdabcdabcdabcdabcdabcdabcdabcd"

		const encrypted = $security.encrypt("123", key)

		const decrypted = $security.decrypt(encrypted, key)

		if (decrypted != "123") {
			throw new Error("Expected decrypted '123', got " + decrypted)
		}
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilesystemBinds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "test")
	}))
	defer srv.Close()

	vm := goja.New()
	vm.Set("mh", &multipart.FileHeader{Filename: "test"})
	vm.Set("testFile", filepath.Join(app.DataDir(), "data.db"))
	vm.Set("baseUrl", srv.URL)
	baseBinds(vm)
	filesystemBinds(vm)

	testBindsCount(vm, "$filesystem", 4, t)

	// fileFromPath
	{
		v, err := vm.RunString(`$filesystem.fileFromPath(testFile)`)
		if err != nil {
			t.Fatal(err)
		}

		file, _ := v.Export().(*filesystem.File)

		if file == nil || file.OriginalName != "data.db" {
			t.Fatalf("[fileFromPath] Expected file with name %q, got %v", file.OriginalName, file)
		}
	}

	// fileFromBytes
	{
		v, err := vm.RunString(`$filesystem.fileFromBytes([1, 2, 3], "test")`)
		if err != nil {
			t.Fatal(err)
		}

		file, _ := v.Export().(*filesystem.File)

		if file == nil || file.OriginalName != "test" {
			t.Fatalf("[fileFromBytes] Expected file with name %q, got %v", file.OriginalName, file)
		}
	}

	// fileFromMultipart
	{
		v, err := vm.RunString(`$filesystem.fileFromMultipart(mh)`)
		if err != nil {
			t.Fatal(err)
		}

		file, _ := v.Export().(*filesystem.File)

		if file == nil || file.OriginalName != "test" {
			t.Fatalf("[fileFromMultipart] Expected file with name %q, got %v", file.OriginalName, file)
		}
	}

	// fileFromUrl (success)
	{
		v, err := vm.RunString(`$filesystem.fileFromUrl(baseUrl + "/test")`)
		if err != nil {
			t.Fatal(err)
		}

		file, _ := v.Export().(*filesystem.File)

		if file == nil || file.OriginalName != "test" {
			t.Fatalf("[fileFromUrl] Expected file with name %q, got %v", file.OriginalName, file)
		}
	}

	// fileFromUrl (failure)
	{
		_, err := vm.RunString(`$filesystem.fileFromUrl(baseUrl + "/error")`)
		if err == nil {
			t.Fatal("Expected url fetch error")
		}
	}
}

func TestFormsBinds(t *testing.T) {
	vm := goja.New()
	formsBinds(vm)

	testBindsCount(vm, "this", 20, t)
}

func TestApisBindsCount(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	apisBinds(vm)

	testBindsCount(vm, "this", 6, t)
	testBindsCount(vm, "$apis", 14, t)
}

func TestApisBindsApiError(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	apisBinds(vm)

	scenarios := []struct {
		js            string
		expectCode    int
		expectMessage string
		expectData    string
	}{
		{"new ApiError()", 0, "", "null"},
		{"new ApiError(100, 'test', {'test': 1})", 100, "Test.", `{"test":1}`},
		{"new NotFoundError()", 404, "The requested resource wasn't found.", "null"},
		{"new NotFoundError('test', {'test': 1})", 404, "Test.", `{"test":1}`},
		{"new BadRequestError()", 400, "Something went wrong while processing your request.", "null"},
		{"new BadRequestError('test', {'test': 1})", 400, "Test.", `{"test":1}`},
		{"new ForbiddenError()", 403, "You are not allowed to perform this request.", "null"},
		{"new ForbiddenError('test', {'test': 1})", 403, "Test.", `{"test":1}`},
		{"new UnauthorizedError()", 401, "Missing or invalid authentication token.", "null"},
		{"new UnauthorizedError('test', {'test': 1})", 401, "Test.", `{"test":1}`},
	}

	for _, s := range scenarios {
		v, err := vm.RunString(s.js)
		if err != nil {
			t.Errorf("[%s] %v", s.js, err)
			continue
		}

		apiErr, ok := v.Export().(*apis.ApiError)
		if !ok {
			t.Errorf("[%s] Expected ApiError, got %v", s.js, v)
			continue
		}

		if apiErr.Code != s.expectCode {
			t.Errorf("[%s] Expected Code %d, got %d", s.js, s.expectCode, apiErr.Code)
		}

		if apiErr.Message != s.expectMessage {
			t.Errorf("[%s] Expected Message %q, got %q", s.js, s.expectMessage, apiErr.Message)
		}

		dataRaw, _ := json.Marshal(apiErr.RawData())
		if string(dataRaw) != s.expectData {
			t.Errorf("[%s] Expected Data %q, got %q", s.js, s.expectData, dataRaw)
		}
	}
}

func TestLoadingDynamicModel(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	baseBinds(vm)
	dbxBinds(vm)
	vm.Set("$app", app)

	_, err := vm.RunString(`
		let result = new DynamicModel({
			text:        "",
			bool:        false,
			number:      0,
			select_many: [],
			json:        [],
			// custom map-like field
			obj: {},
		})

		$app.dao().db()
			.select("text", "bool", "number", "select_many", "json", "('{\"test\": 1}') as obj")
			.from("demo1")
			.where($dbx.hashExp({"id": "84nmscqy84lsi1t"}))
			.limit(1)
			.one(result)

		if (result.text != "test") {
			throw new Error('Expected text "test", got ' + result.text);
		}

		if (result.bool != true) {
			throw new Error('Expected bool true, got ' + result.bool);
		}

		if (result.number != 123456) {
			throw new Error('Expected number 123456, got ' + result.number);
		}

		if (result.select_many.length != 2 || result.select_many[0] != "optionB" || result.select_many[1] != "optionC") {
			throw new Error('Expected select_many ["optionB", "optionC"], got ' + result.select_many);
		}

		if (result.json.length != 3 || result.json[0] != 1 || result.json[1] != 2 || result.json[2] != 3) {
			throw new Error('Expected json [1, 2, 3], got ' + result.json);
		}

		if (result.obj.get("test") != 1) {
			throw new Error('Expected obj.get("test") 1, got ' + JSON.stringify(result.obj));
		}
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadingArrayOf(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	baseBinds(vm)
	dbxBinds(vm)
	vm.Set("$app", app)

	_, err := vm.RunString(`
		let result = arrayOf(new DynamicModel({
			id:   "",
			text: "",
		}))

		$app.dao().db()
			.select("id", "text")
			.from("demo1")
			.where($dbx.exp("id='84nmscqy84lsi1t' OR id='al1h9ijdeojtsjy'"))
			.limit(2)
			.orderBy("text ASC")
			.all(result)

		if (result.length != 2) {
			throw new Error('Expected 2 list items, got ' + result.length);
		}

		if (result[0].id != "84nmscqy84lsi1t") {
			throw new Error('Expected 0.id "84nmscqy84lsi1t", got ' + result[0].id);
		}
		if (result[0].text != "test") {
			throw new Error('Expected 0.text "test", got ' + result[0].text);
		}

		if (result[1].id != "al1h9ijdeojtsjy") {
			throw new Error('Expected 1.id "al1h9ijdeojtsjy", got ' + result[1].id);
		}
		if (result[1].text != "test2") {
			throw new Error('Expected 1.text "test2", got ' + result[1].text);
		}
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpClientBindsCount(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	httpClientBinds(vm)

	testBindsCount(vm, "this", 2, t) // + FormData
	testBindsCount(vm, "$http", 1, t)
}

func TestHttpClientBindsSend(t *testing.T) {
	// start a test server
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("testError") != "" {
			res.WriteHeader(400)
			return
		}

		timeoutStr := req.URL.Query().Get("testTimeout")
		timeout, _ := strconv.Atoi(timeoutStr)
		if timeout > 0 {
			time.Sleep(time.Duration(timeout) * time.Second)
		}

		bodyRaw, _ := io.ReadAll(req.Body)
		defer req.Body.Close()

		// normalize headers
		headers := make(map[string]string, len(req.Header))
		for k, v := range req.Header {
			if len(v) > 0 {
				headers[strings.ToLower(strings.ReplaceAll(k, "-", "_"))] = v[0]
			}
		}

		info := map[string]any{
			"method":  req.Method,
			"headers": headers,
			"body":    string(bodyRaw),
		}

		// add custom headers and cookies
		res.Header().Add("X-Custom", "custom_header")
		res.Header().Add("Set-Cookie", "sessionId=123456")

		infoRaw, _ := json.Marshal(info)

		// write back the submitted request
		res.Write(infoRaw)
	}))
	defer server.Close()

	vm := goja.New()
	baseBinds(vm)
	httpClientBinds(vm)
	vm.Set("testUrl", server.URL)

	_, err := vm.RunString(`
		function getNestedVal(data, path) {
			let result = data || {};
			let parts  = path.split(".");

			for (const part of parts) {
				if (
					result == null ||
					typeof result !== "object" ||
					typeof result[part] === "undefined"
				) {
					return null;
				}

				result = result[part];
			}

			return result;
		}

		let testTimeout;
		try {
			$http.send({
				url:     testUrl + "?testTimeout=3",
				timeout: 1
			})
		} catch (err) {
			testTimeout = err
		}
		if (!testTimeout) {
			throw new Error("Expected timeout error")
		}

		// error response check
		const test0 = $http.send({
			url: testUrl + "?testError=1",
		})

		// basic fields check
		const test1 = $http.send({
			method:  "post",
			url:     testUrl,
			headers: {"header1": "123", "header2": "456"},
			body:    '789',
		})

		// with custom content-type header
		const test2 = $http.send({
			url: testUrl,
			headers: {"content-type": "text/plain"},
		})

		// with FormData
		const formData = new FormData()
		formData.append("title", "123")
		const test3 = $http.send({
			url: testUrl,
			body: formData,
			headers: {"content-type": "text/plain"}, // should be ignored
		})

		const scenarios = [
			[test0, {
				"statusCode": "400",
			}],
			[test1, {
				"statusCode":                "200",
				"headers.X-Custom.0":        "custom_header",
				"cookies.sessionId.value":   "123456",
				"json.method":               "POST",
				"json.headers.header1":      "123",
				"json.headers.header2":      "456",
				"json.headers.content_type": "application/json", // default
				"json.body":                 "789",
			}],
			[test2, {
				"statusCode":                "200",
				"headers.X-Custom.0":        "custom_header",
				"cookies.sessionId.value":   "123456",
				"json.method":               "GET",
				"json.headers.content_type": "text/plain",
			}],
			[test3, {
				"statusCode":                "200",
				"headers.X-Custom.0":        "custom_header",
				"cookies.sessionId.value":   "123456",
				"json.method":               "GET",
				"json.body": [
					"\r\nContent-Disposition: form-data; name=\"title\"\r\n\r\n123\r\n--",
				],
				"json.headers.content_type": [
					"multipart/form-data; boundary="
				],
			}],
		]

		for (let scenario of scenarios) {
			const result = scenario[0];
			const expectations = scenario[1];

			for (let key in expectations) {
				const value = getNestedVal(result, key);
				const expectation = expectations[key]
				if (Array.isArray(expectation)) {
					// check for partial match(es)
					for (let exp of expectation) {
						if (!value.includes(exp)) {
							throw new Error('Expected ' + key + ' to contain ' + exp + ', got: ' + result.raw);
						}
					}
				} else {
					// check for direct match
					if (value != expectation) {
						throw new Error('Expected ' + key + ' ' + expectation + ', got: ' + result.raw);
					}
				}
			}
		}
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCronBindsCount(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	cronBinds(app, vm, nil)

	testBindsCount(vm, "this", 2, t)
}

func TestHooksBindsCount(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	hooksBinds(app, vm, nil)

	testBindsCount(vm, "this", 88, t)
}

func TestHooksBinds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	result := &struct {
		Called int
	}{}

	vmFactory := func() *goja.Runtime {
		vm := goja.New()
		baseBinds(vm)
		vm.Set("$app", app)
		vm.Set("result", result)
		return vm
	}

	pool := newPool(1, vmFactory)

	vm := vmFactory()
	hooksBinds(app, vm, pool)

	_, err := vm.RunString(`
		onModelBeforeUpdate((e) => {
			result.called++;
		}, "demo1")

		onModelBeforeUpdate((e) => {
			throw new Error("example");
		}, "demo1")

		onModelBeforeUpdate((e) => {
			result.called++;
		}, "demo2")

		onModelBeforeUpdate((e) => {
			result.called++;
		}, "demo2")

		onModelBeforeUpdate((e) => {
			return false
		}, "demo2")

		onModelBeforeUpdate((e) => {
			result.called++;
		}, "demo2")

		onAfterBootstrap(() => {
			// check hooks propagation and tags filtering
			const recordA = $app.dao().findFirstRecordByFilter("demo2", "1=1")
			recordA.set("title", "update")
			$app.dao().saveRecord(recordA)
			if (result.called != 2) {
				throw new Error("Expected result.called to be 2, got " + result.called)
			}

			// reset
			result.called = 0;

			// check error handling
			let hasErr = false
			try {
				const recordB = $app.dao().findFirstRecordByFilter("demo1", "1=1")
				recordB.set("text", "update")
				$app.dao().saveRecord(recordB)
			} catch (err) {
				hasErr = true
			}
			if (!hasErr) {
				throw new Error("Expected an error to be thrown")
			}
			if (result.called != 1) {
				throw new Error("Expected result.called to be 1, got " + result.called)
			}
		})

		$app.bootstrap();
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRouterBindsCount(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	routerBinds(app, vm, nil)

	testBindsCount(vm, "this", 3, t)
}

func TestRouterBinds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	result := &struct {
		AddCount int
		UseCount int
		PreCount int
	}{}

	vmFactory := func() *goja.Runtime {
		vm := goja.New()
		baseBinds(vm)
		vm.Set("$app", app)
		vm.Set("result", result)
		return vm
	}

	pool := newPool(1, vmFactory)

	vm := vmFactory()
	routerBinds(app, vm, pool)

	_, err := vm.RunString(`
		routerAdd("GET", "/test", (e) => {
			result.addCount++;
		}, (next) => {
			return (c) => {
				result.addCount++;

				return next(c);
			}
		})

		routerUse((next) => {
			return (c) => {
				result.useCount++;

				return next(c)
			}
		})

		routerPre((next) => {
			return (c) => {
				result.preCount++;

				return next(c)
			}
		})
	`)
	if err != nil {
		t.Fatal(err)
	}

	e, err := apis.InitApi(app)
	if err != nil {
		t.Fatal(err)
	}

	serveEvent := &core.ServeEvent{
		App:    app,
		Router: e,
	}
	if err := app.OnBeforeServe().Trigger(serveEvent); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	e.ServeHTTP(rec, req)

	if result.AddCount != 2 {
		t.Fatalf("Expected AddCount %d, got %d", 2, result.AddCount)
	}

	if result.UseCount != 1 {
		t.Fatalf("Expected UseCount %d, got %d", 1, result.UseCount)
	}

	if result.PreCount != 1 {
		t.Fatalf("Expected PreCount %d, got %d", 1, result.PreCount)
	}
}

func TestFilepathBindsCount(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	filepathBinds(vm)

	testBindsCount(vm, "$filepath", 15, t)
}

func TestOsBindsCount(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := goja.New()
	osBinds(vm)

	testBindsCount(vm, "$os", 17, t)
}
