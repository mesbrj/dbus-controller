
Run as Systemd Service

Create `/etc/systemd/system/dbus-controller.service`:

```ini
[Unit]
Description=D-Bus Controller API
After=network.target

[Service]
Type=simple
User=dbus-controller
ExecStart=/usr/local/bin/dbus-controller
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```