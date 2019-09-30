# weather

```
$ crontab -l
30 * * * * hongkliu /home/hongkliu/bin/weather/weather --debug-mode=true -config=/home/hongkliu/bin/weather/config.yaml > /home/hongkliu/Download/weather_output/weather.log 2>&1
```