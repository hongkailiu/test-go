# weather

```
$ mkdir -p /home/hongkliu/Downloads/weather_output
$ cat /home/hongkliu/bin/weather/config.yaml
appID: <secret>
cities:
  - name: vancouver
    country: ca
  - name: toronto
    country: ca
  - name: montreal
    country: ca
writer:
  - logger
  - csv
  - yaml
outputDir: /home/hongkliu/Downloads/weather_output

### hourly
$ crontab -l
0 * * * * /home/hongkliu/bin/weather/weather --debug-mode=true -config=/home/hongkliu/bin/weather/config.yaml >> /home/hongkliu/Downloads/weather_output/weather.log 2>&1

```