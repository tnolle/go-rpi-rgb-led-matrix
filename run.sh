#!/bin/bash

# Default settings
MODE=0
ROWS=32
COLS=32
BRIGHTNESS=40
GPIO_MAPPING=adafruit-hat
MULTIPLEXING=0
LED_CHAIN=1
ROW_ADDR_TYPE=3


function run_demo() {
  sudo ./3rdparty/rpi-rgb-led-matrix/examples-api-use/demo \
    -D$MODE \
    --led-gpio-mapping=$GPIO_MAPPING \
    --led-cols=$COLS \
    --led-rows=$ROWS \
    --led-chain=$LED_CHAIN \
    --led-multiplexing=$MULTIPLEXING \
    --led-row-addr-type=$ROW_ADDR_TYPE \
    --led-brightness=$BRIGHTNESS \
    --led-slowdown-gpio=4
}

run_demo