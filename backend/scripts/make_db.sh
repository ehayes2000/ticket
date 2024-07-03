#! /bin/bash
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)"
DB_FILE=$HERE/../db.db
SCHEMA=$HERE/schema.sql
sqlite3 $DB_FILE < $SCHEMA



