package main

import (
	"machine"
	"time"
)

func main() {
	// Configure ADC0 and ADC1 as analog inputs
	adcA0 := machine.ADC{Pin: machine.ADC0}
	adcA1 := machine.ADC{Pin: machine.ADC1}

	// Initialize the ADCs
	adcA0.Configure(machine.ADCConfig{})
	adcA1.Configure(machine.ADCConfig{})

	// Initialize the UART connection
	uart := machine.UART0
	uart.Configure(machine.UARTConfig{})

	for {
		// Read voltage on ADC0
		voltageA0 := adcA0.Get() // Assuming 10-bit ADC

		// Read voltage on ADC1
		voltageA1 := adcA1.Get() // Assuming 10-bit ADC

		// Write data to UART
		uart.WriteByte(byte(voltageA0 >> 8))
		uart.WriteByte(byte(voltageA0 & 0xFF))
		uart.WriteByte('\n')
		uart.WriteByte(byte(voltageA1 >> 8))
		uart.WriteByte(byte(voltageA1 & 0xFF))
		uart.WriteByte('\n')

		time.Sleep(time.Second)
	}
}
