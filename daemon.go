package aged

import (
	"log"
	"os"
	"syscall"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

// Daemon is the main aged daemon type.
type Daemon struct {
	conn        *dbus.Conn
	running     bool
	AgeProvider AgeProvider
}

// NewDaemon creates a new daemon, returns error if it cannot connect to the dbus session bus.
func NewDaemon() (*Daemon, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}
	return &Daemon{conn: conn}, nil
}

// Running returns true if the daemon is in a running state
func (d *Daemon) Running() bool {
	return d.running
}

func (d *Daemon) shutdown() {
	d.running = false
	d.conn.Close()
}

func (d *Daemon) reload() {

}

// HandleSignal handles os signal
func (d *Daemon) HandleSignal(sig os.Signal) {
	if sig == syscall.SIGHUP {
		d.reload()
		return
	}
	if sig == os.Interrupt {
		d.shutdown()
		return
	}
}

// Iface is the dbus interface
const Iface = "org.freedesktop.AgeVerification1"

// Path is the dbus path
const Path = "/org/freedesktop/AgeVerification1"

// Schema is the dbus XML definition
const Schema = `<node>
	<interface name="org.freedesktop.AgeVerification1">
		<method name="GetAgeBracket">
			<arg direction="out" type="s"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `</node> `

// SetupHandler sets up internal dbus handler.
func (d *Daemon) SetupHandler() {
	log.Println("Setting up DBus Handler")

	err := d.conn.ExportMethodTable(map[string]any{
		"GetAgeBracket": d.GetAgeBracket,
	}, Path, Iface)

	if err != nil {
		log.Fatalf("Failed to export DBus interface 1: %s", err.Error())
	}
	err = d.conn.Export(introspect.Introspectable(Schema), Path, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		log.Fatalf("Failed to export DBus interface 2: %s", err.Error())
	}
	reply, err := d.conn.RequestName(Iface, dbus.NameFlagDoNotQueue)
	if reply != dbus.RequestNameReplyPrimaryOwner {
		if err != nil {
			log.Println(err.Error())
		}
		os.Exit(1)
	}

	log.Println("Running")

	d.running = true
}

// Under13 denotes the user is under 13 years of age.
const Under13 = "<13"

// Between13And16 denotes the user is between 13 and 16 years of age.
const Between13And16 = ">=13,<=16"

// Seventeen denotes the user is 17 years old.
const Seventeen = "=17"

// Adult denotes that the user is above 18 years of age.
const Adult = ">=18"

// USAAdult denotes that the user is above 21 years of age.
const USAAdult = ">=21"

// GetAgeBracket handles a dbus method call to determine the age of the user who owns the session bus.
func (d *Daemon) GetAgeBracket() (string, *dbus.Error) {
	age, err := d.AgeProvider.LookupUserAge()
	if err != nil {
		return "", err
	}
	bracket := USAAdult
	if age < 21 {
		bracket = Adult
	}
	if age < 18 {
		bracket = Seventeen
	}
	if age <= 16 {
		bracket = Between13And16
	}
	if age < 13 {
		bracket = Under13
	}
	return bracket, nil
}
