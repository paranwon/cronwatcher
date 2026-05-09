# cronwatcher

Lightweight daemon that monitors cron job execution, logs durations, and alerts on missed or long-running jobs.

---

## Installation

```bash
go install github.com/youruser/cronwatcher@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/cronwatcher.git && cd cronwatcher && go build -o cronwatcher .
```

---

## Usage

Start the daemon with a config file:

```bash
cronwatcher --config /etc/cronwatcher/config.yaml
```

Example `config.yaml`:

```yaml
jobs:
  - name: daily-backup
    schedule: "0 2 * * *"
    timeout: 30m
    alert_on_miss: true

  - name: hourly-sync
    schedule: "0 * * * *"
    timeout: 5m
    alert_on_miss: true

alerts:
  webhook: "https://hooks.example.com/notify"

log:
  path: /var/log/cronwatcher.log
  level: info
```

Wrap your existing cron commands to report execution status:

```bash
# In your crontab
0 2 * * * cronwatcher exec --job daily-backup -- /usr/local/bin/backup.sh
```

cronwatcher will log start time, duration, and exit code, and fire alerts if the job exceeds its timeout or fails to run on schedule.

---

## License

MIT © youruser