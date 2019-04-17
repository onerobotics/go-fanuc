# go-fanuc

go-fanuc is a Go client library for accessing FANUC robots over HTTP.
The plan is to eventually support alternative clients (e.g. a filesystem
client).

## Usage

    import "github.com/onerobotics/go-fanuc/fanuc"

Construct a new FANUC client, then use the various services on the client
to access different parts of the FANUC robot. For example:

    c := &fanuc.Client{Host: "127.0.0.101"}

    numregs, err := c.GetNumericRegisters()
    if err != nil {
    	panic(err)
    }

    for _, r := range numregs {
    	fmt.Println(r.String())
    }
