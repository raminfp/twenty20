# Twenty20 Service (Linux)

## Prerequisites
- Go programming language installed
- Sudo access on the target system
- Basic understanding of systemd service management

## Deployment Steps

### 1. Build the Service
Compile the Go application:
```bash
go build -o twenty20_service main.go
```

### 2. Install the Binary
Move the compiled binary to the system-wide bin directory:
```bash
sudo mv twenty20_service /usr/local/bin/
```

### 3. Create Systemd Service File
Create a systemd service configuration:
```bash
sudo nano /etc/systemd/system/twenty20.service
```

Paste the following configuration:
```ini
[Unit]
Description=twenty20_service
After=network.target

[Service]
ExecStart=/usr/local/bin/twenty20_service
Restart=always
User=raminfp
Group=raminfp
Environment=DISPLAY=:0
Environment=DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus
WorkingDirectory=/home/raminfp

[Install]
WantedBy=multi-user.target
```

### 4. Configure Systemd
Reload systemd configuration and enable the service:
```bash
# Reload systemd daemon
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable twenty20.service

# Start the service
sudo systemctl start twenty20.service

# Status the service
sudo systemctl status twenty20.service
```

### 5. Verify Service Status
Check the service is running correctly:
```bash
sudo systemctl status twenty20.service
```

## Service Management Commands
- Start service: `sudo systemctl start twenty20.service`
- Stop service: `sudo systemctl stop twenty20.service`
- Restart service: `sudo systemctl restart twenty20.service`
- Check status: `sudo systemctl status twenty20.service`

## Troubleshooting
- Ensure the binary is executable
- Check system logs with `journalctl -u twenty20.service`
- Verify user and group permissions

## Notes
- Modify the `User` and `Group` fields if needed
- Adjust `After=` directive based on service dependencies
