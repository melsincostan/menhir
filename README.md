# Menhir

Work in progress!

Sometimes developing a website implies doing `window.fetch()` calls against legacy APIs on a different domain doesn't work due to them not setting CORS headers.
Instead of adjusting a productive application, this offers a way to add the headers by means of a reverse proxy.
It can do more than this, and writing new modules for it should (hopefully) be quite easy by implementing one (or multiple) of the interfaces in `menhir/interfaces.go`.

Runs on `0.0.0.0:8080` by default, with only a logging module loaded (every request passed to the backend unchanged - besides setting the "host" header), with `http://localhost` as a backend.
Set another destination host using the `--destination` flag, a listen port using `--port` or a bind host using `host`.
Modules can also be enabled / disbled / configured over flags. Use `--help` for a list of all flags, or `--list-modules` to list all modules.
