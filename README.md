# Datadis Go Client

**⚠️ Work in progress, experimental**

Go API to read [Datadis](https://datadis.es) energy consumption info.

You'll need https://www.datadis.es account to use this package.

## Example

```Go
package main

import (
        "fmt"
        "os"
        "time"

        "github.com/rubiojr/go-datadis"
)

// Fetch datadis last day consumption
func main() {
        client := datadis.NewClient()
        client.Login(os.Getenv("DATADIS_USERNAME"), os.Getenv("DATADIS_PASSWORD"))
        s, err := client.Supplies()
        if err != nil {
                panic(err)
        }

        now := time.Now()
        year, month, day := now.Date()
        // Read yesterday's data
        date := time.Date(year, month, day-1, 0, 0, 0, 0, now.UTC().Location())
        data, err := client.ConsumptionData(&s[0], date, date)
        for _, d := range data {
                fmt.Println("CUPS: ", d.Cups)
                fmt.Println("Date: ", d.Date)
                fmt.Println("Time: ", d.Time)
                fmt.Printf("Consumption: %f KWh\n", d.Consumption)
                fmt.Println("Obtained Method: ", d.ObtainMethod)
        }
}
```

## Building the command line client

```
make
```

## Using the sample client

Export the datadis username and password as environment variables:

```
export DATADIS_USERNAME="username here"
export DATADIS_PASSWORD="password here"
```

Run the client without arguments:

```
./bin/datadis
```

## Related

* https://github.com/uvejota/homeassistant-edata
* https://github.com/rubiojr/go-edistribucion
* https://github.com/azogue/aiopvpc
* https://github.com/trocotronic/edistribucion
* https://github.com/uvejota/edistribucion
