#!/bin/bash

docker run -d --name mailhog -p 8025:8025 -p 1025:1025 mailhog/mailhog > /dev/null
echo Press enter to run the watch.
read

elastic-watcher --config ./config.yml run --now ./watch.yml

if [ $? == 0 ]; then
  echo Successfully sent email.  You can see it from here:
  echo http://localhost:8025/
  echo Press enter to exit
  read
else
  echo Failed to send email.
fi

docker rm -f mailhog > /dev/null