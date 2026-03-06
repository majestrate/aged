# AgeD

DBus Based Age Verification Toy Program.

Runs on the session bus for each user.

## Building

```bash
$ go get github.com/majestrate/aged
```

## Usage

First populate the date of birth for the current user:

    $ echo "2001-09-11" > ~/.config/birthday

Run the daemon:

    $ aged

Call via `dbus-send`:

    $ dbus-send --print-reply --session --dest=org.freedesktop.AgeVerification1 /org/freedesktop/AgeVerification1 org.freedesktop.AgeVerification1.GetAgeBracket 

