#!/bin/bash

if [ ! -d /etc/systemd/system ] || [ ! -d /usr/bin ] || [ ! -d /usr/lib ]; then
  echo "Your system is not eligible for this install script, please install manually"
fi

if [ "$EUID" -ne 0 ];  then
  echo "Please run as root"
  exit
fi

exec 2> /dev/null

echo -e "\n.: WELCOME TO MICROHTTP :."

echo -e "\nCleaning..."
rm -f /etc/systemd/system/microhttp.service
rm -f /usr/bin/microhttp

echo "Checking for micro-www user"
if [ "$(id -u micro-www)" == "" ]; then
  echo ">>>Making user"
  useradd -d /usr/lib/microhttp -m -s /sbin/nologin micro-www
fi

echo "Checking required files"
if [ ! -e /sbin/setcap ]; then
  echo "Make sure setcap is installed"
  exit
fi
if ! hash unzip 2> /dev/null || ! hash unzip 2> /dev/null; then
  echo "Make sure unzip and wget are installed"
  exit
fi


echo "Installing"
cp ./opt/microhttp_linux64 /usr/bin/microhttp
cp ./opt/systemd/microhttp.service /etc/systemd/system/microhttp.service
setcap cap_net_bind_service=+ep /usr/bin/microhttp
mkdir -p /usr/lib/microhttp/www

if [ ! -e /usr/lib/microhttp/main.json ]; then
  cp ./opt/config/example.json /usr/lib/microhttp/main.json
else
  cp -f ./opt/config/example.json /usr/lib/microhttp/main.json.new
fi

if [ ! -e /usr/lib/microhttp/www/index.html ]; then
  echo "<html><head></head><body><h1>Welcome to MicroHTTP!</h1></body></html>" > /usr/lib/microhttp/www/index.html
fi

systemctl stop microhttp.service
systemctl enable microhttp.service
systemctl start microhttp.service

echo -e "\nAnd we are done. Make sure to check the configuration at /usr/lib/microhttp/main.json\n"
