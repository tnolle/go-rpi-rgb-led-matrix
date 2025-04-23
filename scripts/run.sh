#!/bin/bash

# Default settings
MODE=0
ROWS=32
COLS=64
BRIGHTNESS=33
GPIO_MAPPING=regular
MULTIPLEXING=2
LED_CHAIN=3
PARALLEL=3
ROW_ADDR_TYPE=3


function demo() {
  sudo ./3rdparty/rpi-rgb-led-matrix/examples-api-use/demo \
    -D$MODE \
    --led-gpio-mapping=$GPIO_MAPPING \
    --led-cols=$COLS \
    --led-rows=$ROWS \
    --led-parallel=$PARALLEL \
    --led-chain=$LED_CHAIN \
    --led-multiplexing=$MULTIPLEXING \
    --led-row-addr-type=$ROW_ADDR_TYPE \
    --led-brightness=$BRIGHTNESS \
    --led-slowdown-gpio=-1
}

function autodarts() {
  sudo ./3rdparty/rpi-rgb-led-matrix/examples-api-use/demo \
    -D1 -m0 images/out/autodarts.ppm \
    --led-gpio-mapping=$GPIO_MAPPING \
    --led-cols=$COLS \
    --led-rows=$ROWS \
    --led-parallel=$PARALLEL \
    --led-chain=$LED_CHAIN \
    --led-multiplexing=$MULTIPLEXING \
    --led-row-addr-type=$ROW_ADDR_TYPE \
    --led-brightness=$BRIGHTNESS \
    --led-show-refresh \
    --led-limit-refresh=120 \
    --led-slowdown-gpio=4
}

function celonis() {
 sudo ./3rdparty/rpi-rgb-led-matrix/examples-api-use/demo \
   -D1 -m0 images/out/celonis.ppm \
   --led-gpio-mapping=$GPIO_MAPPING \
   --led-cols=$COLS \
   --led-rows=$ROWS \
   --led-parallel=$PARALLEL \
   --led-chain=$LED_CHAIN \
   --led-multiplexing=$MULTIPLEXING \
   --led-row-addr-type=$ROW_ADDR_TYPE \
   --led-brightness=$BRIGHTNESS \
   --led-show-refresh \
   --led-limit-refresh=120 \
   --led-slowdown-gpio=4
}

function text() {
  sudo ./3rdparty/rpi-rgb-led-matrix/examples-api-use/text-example -f ./3rdparty/rpi-rgb-led-matrix/fonts/5x7.bdf \
      --led-gpio-mapping=$GPIO_MAPPING \
      --led-cols=$COLS \
      --led-rows=$ROWS \
      --led-parallel=$PARALLEL \
      --led-chain=$LED_CHAIN \
      --led-multiplexing=$MULTIPLEXING \
      --led-row-addr-type=$ROW_ADDR_TYPE \
      --led-brightness=$BRIGHTNESS \
      --led-slowdown-gpio=4
}

text