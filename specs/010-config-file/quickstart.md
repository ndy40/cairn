# Quickstart: Configuration File Support

**Feature**: 010-config-file

## Scenario 1: Config file overrides default

1. Create config directory: `mkdir -p ~/.config/cairn`
2. Create config file: `echo '{"db_path": "/tmp/cairn-test.db"}' > ~/.config/cairn/cairn.json`
3. Run: `cairn config`
4. **Expected**: `CAIRN_DB_PATH=/tmp/cairn-test.db`
5. Cleanup: `rm ~/.config/cairn/cairn.json`

## Scenario 2: Env var overrides config file

1. Create config file: `echo '{"db_path": "/tmp/from-file.db"}' > ~/.config/cairn/cairn.json`
2. Run: `CAIRN_DB_PATH=/tmp/from-env.db cairn config`
3. **Expected**: `CAIRN_DB_PATH=/tmp/from-env.db` (env var wins)
4. Cleanup: `rm ~/.config/cairn/cairn.json`

## Scenario 3: CLI flag overrides config file (but not env var)

1. Create config file: `echo '{"db_path": "/tmp/from-file.db"}' > ~/.config/cairn/cairn.json`
2. Run: `cairn --db /tmp/from-flag.db config`
3. **Expected**: `CAIRN_DB_PATH=/tmp/from-flag.db` (CLI flag wins over config file)
4. Run: `CAIRN_DB_PATH=/tmp/from-env.db cairn --db /tmp/from-flag.db config`
5. **Expected**: `CAIRN_DB_PATH=/tmp/from-env.db` (env var wins over CLI flag)
6. Cleanup: `rm ~/.config/cairn/cairn.json`

## Scenario 4: No config file — works as before

1. Ensure no config file: `rm -f ~/.config/cairn/cairn.json`
2. Run: `cairn config`
3. **Expected**: `CAIRN_DB_PATH=` followed by the OS default path (e.g., `~/.local/share/cairn/bookmarks.db`)

## Scenario 5: Invalid config file

1. Create invalid config: `echo 'not json' > ~/.config/cairn/cairn.json`
2. Run: `cairn config`
3. **Expected**: Error message on stderr mentioning invalid JSON, exit code 1
4. Cleanup: `rm ~/.config/cairn/cairn.json`

## Scenario 6: Dropbox app key from config file

1. Create config file: `echo '{"dropbox_app_key": "test-key-123"}' > ~/.config/cairn/cairn.json`
2. Run: `cairn config`
3. **Expected**: `CAIRN_DROPBOX_APP_KEY=test-key-123`
4. Cleanup: `rm ~/.config/cairn/cairn.json`
