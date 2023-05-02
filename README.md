# tmenu

A frontend for a text-mode terminal server.

## Command-line Usage

Customize "banner.txt" with whatever the connection banner should say,
and customize "targets.json" with what you want to connect to.

Now you can just do `./tmenu -targetfile /path/to/targets.json
-bannerfile /path/to/banner.txt` and then use the menu to connect.

Granted, this is only slightly more convenient than connecting to a
port.

However, now you could have a guest user that had tmenu as its shell,
if that was a thing you wanted to do.

## Web Access Usage

If you use `ttyd` (https://github.com/tsl0922/ttyd) or `gotty`
(https://github.com/yudai/gotty) you can make access to the menu
available via the web.  I recommend `ttyd` because it is actively maintained.

### ttyd

Use the `ttyd.service` file here as a systemd unit to hook up gotty
on port 6180 to tmenu (assuming you've installed tmenu in
/usr/local/bin and its support files in /usr/local/share/tmenu).

### gotty

You will probably need https://github.com/yudai/gotty/pull/259 to make
gotty compile on a modern Go.

Use the `gotty.service` file here as a systemd unit to hook up gotty
on port 6180 to tmenu (assuming you've installed tmenu in
/usr/local/bin and its support files in /usr/local/share/tmenu).

### Exposing the terminal-to-HTTP interface via a web server

Now port 6180 serves the menu, so pointing a browser at
http://localhost:6180 should work.  If you're using something other
than systemd, or just hate systemd (and who could blame you?), you
basically just want to set tmenu up as an inetd service.

If you want to let other people use this you probably want to put that
on a route behind TLS and on port 443 of a real web server.  I use
Apache, so if you use Nginx or something else, you'll need to figure
out the way to expose it on your own.  Fundamentally all you're doing
is setting up a reverse proxy that supports websockets.

For Apache 2.4, you need `mod_proxy` and `mod_proxy_wstunnel` and will
want something like this inside your virtualhost definition (you will
also want all the standard stuff to turn on TLS, and expose
.well-known if you're using Let's Encrypt, and all that jazz--this is
just the extra config on top of that to turn on proxying):

```
	ProxyVia On
	ProxyPreserveHost On
    	ProxyRequests off
	RewriteEngine on
	RewriteCond %{HTTP:Upgrade} websocket [NC,OR]
	RewriteCond %{HTTP:Connection} upgrade [NC,OR]

	<Location />
		  Order allow,deny
		  Allow from all
	</Location>

	ProxyPass /ws ws://localhost:6180/ws
	ProxyPassReverse /ws ws:localhost:6180/ws
	ProxyPass / http://localhost:6180/
	ProxyPassReverse / http://localhost:6180/
```

Good luck!
