# Aweeing
Приложенька для синка ical -> [awtrix-light](https://blueforcer.github.io/awtrix-light/).
Использую в личных нуждах, чтоб показать домочадцам, когда не стоит выбивать дверь кабинета.

## Примерчик
На примере моего домашнего:
```yaml
verbose: false
calendar:
  sourceUrl: "https://calendar.yandex.ru/export/ics.xml?private_token=XXXXX&tz_id=Asia/Bangkok"
  timezone: Asia/Bangkok
ticker:
  jitter: 20m
  previewLimit: 24h
  fetchInterval: 1h
  tickInterval: 5m
mqtt:
  upstream: "tcp://mqtt.iot.buglloc.cc:1883"
  username: "aweeting"
  password: "aweeting-password"
  topic: "awtrix/custom/meetings"
awtrix:
  selfDestruct: true
  upcomingLimit: 4h
  messages:
    none:
      color: "#ffffff"
      icon: "2899"
    upcoming:
      color: "#ffffff"
      icon: "11899"
    onAir:
      color: "#e60000"
      icon: "24092"
```

Будет получен следующий результат:
  - встречи с перерывами менее 20 минут (`ticker.jitter`) будут объеденены в один интервал
  - приложенька само выпиливается (`awtrix.selfDestruct`), если нет запланированных встреч ближайшие 4 часа (`awtrix.upcomingLimit`)
  - для запланированных встреч показываем иконку ["Terminator Eye"](https://developer.lametric.com/content/apps/icon_thumbs/11899_icon_thumb.gif) (`awtrix.messages.upcoming.icon`) и время _до_ встречи белым `awtrix.messages.upcoming.color`)
  - для идущей встречи показываем иконку ["terminator eye glow"](https://developer.lametric.com/content/apps/icon_thumbs/24092_icon_thumb.gif) (`awtrix.messages.onAir.icon`) и время до окончания встречи красненьким `awtrix.messages.onAir.color`)

Примерчики:
  - встреча начнется через 13 минут:
![upcoming.gif](example%2Fupcoming.gif)
  - встреча закончится через час: 
![on-air.gif](example%2Fon-air.gif) 