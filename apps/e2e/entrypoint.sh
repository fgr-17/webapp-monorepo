#!/bin/sh
set -eu

BASE_URL="${BASE_URL:-http://frontend}"
SELENIUM_URL="${SELENIUM_REMOTE_URL:-http://chrome:4444/wd/hub}"

echo "Waiting for frontend at ${BASE_URL} ..."
i=0
until curl -fsS "${BASE_URL}/" >/dev/null 2>&1; do
  i=$((i + 1))
  if [ "$i" -gt 60 ]; then
    echo "Frontend did not become ready in time" >&2
    exit 1
  fi
  sleep 2
done

echo "Waiting for Selenium at ${SELENIUM_URL} ..."
i=0
until curl -fsS "http://chrome:4444/status" >/dev/null 2>&1; do
  i=$((i + 1))
  if [ "$i" -gt 60 ]; then
    echo "Selenium did not become ready in time" >&2
    exit 1
  fi
  sleep 2
done

echo "Running Behave ..."
exec behave "$@"
