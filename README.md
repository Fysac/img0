# img0

img0 is a fast, lightweight, standalone image hosting program with zero external dependencies, inspired by the emergence of complexity and dark patterns in popular image hosts.
 The goal is to do one thing well: host directly-linkable images.

### Building
Simply run `make` or `go build img0.go`.

### Running
For testing, you can just launch the img0 binary. It runs on http://127.0.0.1:8000 by default.

If you want to run img0 in a production setting, I recommend setting it up behind an nginx reverse proxy. This gives you more flexibility in areas like configuring TLS and including custom assets.

The nginx proxy configuration might look something like this:

    location /img0 {
        proxy_buffering off;
		proxy_pass http://127.0.0.1:8000;
    }

### Images

* Uploaded images are stored in `img/`
* EXIF data is stripped by default for privacy; toggle constant `ReencodeJPEG` to disable, at the expense of possible greater file size
* Max size is 10 MB by default; use the constant `MaxRequestSize` to adjust as needed
