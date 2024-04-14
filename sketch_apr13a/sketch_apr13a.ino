#include "LiquidCrystal.h"

LiquidCrystal lcd(8, 9, 4, 5, 6, 7);

void setup() {
   Serial.begin(9600);
   lcd.begin(16, 2);
   lcd.print("The Pertubators!");
}

float readAndConvertAnalog(int pin) {
   int analog_value = analogRead(pin);
   float voltage = (analog_value * 5.0) / 1024.0;
   return voltage;
}
float readAndConvertAnalogCurrent(int pin) {
  int analog_value = analogRead(pin);
   float current = (analog_value * 5.0) / 1024.0;
   return current;
}

void printVoltage(float voltage,float current) {
    Serial.print("v= ");
    Serial.println(voltage, 4); // Print voltage with 4 decimal places to the serial monitor

    lcd.setCursor(0, 1); // Set cursor to the start of the second row
    lcd.print(String(voltage, 4)); // Convert voltage to a string with 4 decimal places and print it
    lcd.print("V ");
    lcd.print(String(current,4));
    lcd.print("mA");
}

void loop() {
   float voltage1 = readAndConvertAnalog(A0);
   float current1 = readAndConvertAnalog(A1);

   printVoltage(voltage1,current1);

   delay(300);
}
