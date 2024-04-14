package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/jacobsa/go-serial/serial"
)

func main() {
	// Set up options for the serial port
	options := serial.OpenOptions{
		PortName:        "/dev/ttyACM0",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the serial port
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Create a new CSV file
	file, err := os.Create("voltage_readings.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Write header to CSV file
	header := []string{"Timestamp", "Voltage on ADC0 (V)", "Voltage on ADC1 (V)"}
	err = writer.Write(header)
	if err != nil {
		log.Fatal(err)
	}
	writer.Flush()

	// Buffer to hold incoming data
	buf := make([]byte, 2)

	for {
		// Read voltage on ADC0
		_, err := port.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}
		voltageA0 := int(buf[0])<<8 | int(buf[1])

		// Read voltage on ADC1
		_, err = port.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}
		voltageA1 := int(buf[0])<<8 | int(buf[1])

		// Write data to CSV file
		record := []string{strconv.Itoa(voltageA0), strconv.Itoa(voltageA1)}
		err = writer.Write(record)
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
	}
}
