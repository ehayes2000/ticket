#!/bin/sh
HERE="$(dirname "$0")"
# HERE="$(cd "$(dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)"
DB_FILE=$HERE/../db.db
SCHEMA=$HERE/schema.sql
SEED=$HERE/seed.sql
sqlite3 $DB_FILE < $SCHEMA
sqlite3 $DB_FILE < $SEED


