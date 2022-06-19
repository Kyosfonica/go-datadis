package storage

import "github.com/rubiojr/go-datadis"

type Measurements []*datadis.Measurement

type Storage interface {
	Reader
	Writer
	Lister
}

type Reader interface {
	Read(id int) (*datadis.Measurement, error)
}

type ReaderFunc func(id int) (*datadis.Measurement, error)

func (f ReaderFunc) Read(id int) (*datadis.Measurement, error) {
	return f(id)
}

type Writer interface {
	Write(measurement datadis.Measurement) error
}

type WriterFunc func(measurement datadis.Measurement) error

func (f WriterFunc) Write(measurement datadis.Measurement) error {
	return f(measurement)
}

func DummyWriter() WriterFunc {
	return func(measurement datadis.Measurement) error {
		return nil
	}
}

type Lister interface {
	List(query string) (Measurements, error)
}

type ListerFunc func(query string) (Measurements, error)

func (f ListerFunc) List(query string) (Measurements, error) {
	return f(query)
}
