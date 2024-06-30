# GoCheckIP server

Basic HTTP server to show visitor IP-address.

It just shows
```Your IP is: <ip address>```

Can be used for revealing IP-address for example when using mobile connections, or getting your dynamic address for updating to DNS.

Supports IPv4 and IPv6 addresses, based on your server config.

## Reverse proxy support
Reverse proxy headers are supported in the format of _X-Real-Ip_ and _X-Forwarded-For_. Note: get IP from headers only if you control the proxy server, otherwise it can be spoofed. This app picks the last "X-Forwarded-For" value from the list and will not really validate the IP-address.

Enable proxy usage by providing `-proxy` parameter.

Example config line for [nginx](https://nginx.org/):
```
location /checkip {
        proxy_set_header   X-Real-IP         $remote_addr;
        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_pass         http://localhost:8080/;
}
```
`X-Real-Ip` will be used first if both headers are present.

## Systemd service
There is an example [systemd](https://systemd.io/) service file. Copy it your systemd folder (for example `/etc/systemd/system/`). Edit the `ExecStart` to point to gocheckip-executable. Then `systemctl daemon-reload`, `systemctl enable gocheckip.service` and finally `systemctl start gocheckip.service`.

## Header logging
All request headers can be logged by providing `-verbose` parameter.

---
Made by Visa Hannula
