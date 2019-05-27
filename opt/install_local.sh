#!/bin/bash

if [ ! -d /etc/systemd/system ] || [ ! -d /usr/bin ] || [ ! -d /usr/lib ]; then
  echo "Your system is not eligible for this install script, please install manually"
fi

if [ "$EUID" -ne 0 ];  then
  echo "Please run as root"
  exit
fi

exec 2> /dev/null

echo -e "\n.: WELCOME TO MAGUROHTTP :."

echo -e "\nCleaning..."
rm -f /etc/systemd/system/magurohttp.service
rm -f /usr/bin/magurohttp

echo "Checking for tuna-www user"
if [ "$(id -u tuna-www)" == "" ]; then
  echo ">>>Making user"
  useradd -d /usr/lib/magurohttp -m -s /sbin/nologin tuna-www
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
cp ./opt/magurohttp_linux64 /usr/bin/magurohttp
cp ./opt/systemd/magurohttp.service /etc/systemd/system/magurohttp.service
setcap cap_net_bind_service=+ep /usr/bin/magurohttp
mkdir -p /usr/lib/magurohttp/www

if [ ! -e /usr/lib/magurohttp/main.json ]; then
  cp ./opt/config/example.json /usr/lib/magurohttp/main.json
else
  cp -f ./opt/config/example.json /usr/lib/magurohttp/main.json.new
fi

if [ ! -e /usr/lib/magurohttp/www/index.html ]; then
  echo "<html><head></head><body><h1>Welcome to MaguroHTTP!</h1></body></html>" > /usr/lib/magurohttp/www/index.html
fi

systemctl stop magurohttp.service
systemctl enable magurohttp.service
systemctl start magurohttp.service

echo -e "\nAnd we are done. Make sure to check the configuration at /usr/lib/magurohttp/main.json\n"
