PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE `_migrations` (file VARCHAR(255) PRIMARY KEY NOT NULL, applied INTEGER NOT NULL);
INSERT INTO _migrations VALUES('1640988000_init.go',1679427415788852);
INSERT INTO _migrations VALUES('1660821103_add_user_ip_column.go',1679427415792425);
INSERT INTO _migrations VALUES('1677760279_uppsercase_method.go',1679427415792706);
INSERT INTO _migrations VALUES('1699187560_logs_generalization.go',1700504854831333);
CREATE TABLE `_logs` (
				`id`      TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL,
				`level`   INTEGER DEFAULT 0 NOT NULL,
				`message` TEXT DEFAULT "" NOT NULL,
				`data`    JSON DEFAULT "{}" NOT NULL,
				`created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
				`updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL
			);
CREATE INDEX _logs_level_idx on `_logs` (`level`);
CREATE INDEX _logs_message_idx on `_logs` (`message`);
CREATE INDEX _logs_data_auth_idx on `_logs` (JSON_EXTRACT(`data`, '$.auth'));
CREATE INDEX _logs_created_hour_idx on `_logs` (strftime('%Y-%m-%d %H:00:00', `created`));
CREATE INDEX idx_logs_level on `_logs` (`level`);
CREATE INDEX idx_logs_message on `_logs` (`message`);
CREATE INDEX idx_logs_created_hour on `_logs` (strftime('%Y-%m-%d %H:00:00', `created`));
COMMIT;
