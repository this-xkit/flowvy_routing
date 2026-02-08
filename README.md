# Flowvy Routing

Оптимизированные правила маршрутизации для **Xray**, **V2Ray** и **Mihomo** (Clash).

Автоматически собирает `geosite.dat`, `geoip.dat` и YAML файлы из проверенных источников. Идеально подходит для iOS-клиентов с ограничениями по памяти.

## Скачать

| Формат | Файлы |
|--------|-------|
| **Xray / V2Ray** | [geosite.dat](https://github.com/this-xkit/flowvy_routing/releases/latest/download/geosite.dat) &#124; [geoip.dat](https://github.com/this-xkit/flowvy_routing/releases/latest/download/geoip.dat) |
| **Mihomo** | [category-ru.yaml](https://github.com/this-xkit/flowvy_routing/releases/latest/download/category-ru.yaml) &#124; [category-ru-ip.yaml](https://github.com/this-xkit/flowvy_routing/releases/latest/download/category-ru-ip.yaml) &#124; [apple.yaml](https://github.com/this-xkit/flowvy_routing/releases/latest/download/apple.yaml) &#124; [apple-ip.yaml](https://github.com/this-xkit/flowvy_routing/releases/latest/download/apple-ip.yaml) |
| **Plaintext** | [category-ru.txt](https://github.com/this-xkit/flowvy_routing/releases/latest/download/category-ru.txt) &#124; [category-ru-ip.txt](https://github.com/this-xkit/flowvy_routing/releases/latest/download/category-ru-ip.txt) |

## Категории

### geosite.dat (домены)

| Категория | Описание |
|-----------|----------|
| `category-ru` | Российские домены (.ru, .su, .рф, .by) + популярные сервисы (Яндекс, VK, Ozon, Wildberries и др.) |
| `apple` | Сервисы Apple |
| `private` | Локальные домены (localhost, local, lan, intranet) |

### geoip.dat (IP-адреса)

| Категория | Описание |
|-----------|----------|
| `category-ru-ip` | IP-адреса России |
| `category-by-ip` | IP-адреса Беларуси |
| `category-kz-ip` | IP-адреса Казахстана |
| `apple-ip` | IP-адреса Apple Inc. (17.0.0.0/8 + IPv6) |
| `private-ip` | Приватные IP (10.x, 192.168.x, 172.16-31.x) |

## Примеры конфигурации

### Xray / V2Ray

```json
{
  "routing": {
    "domainStrategy": "IPIfNonMatch",
    "rules": [
      {
        "type": "field",
        "domain": ["geosite:private"],
        "outboundTag": "direct"
      },
      {
        "type": "field",
        "ip": ["geoip:private-ip"],
        "outboundTag": "direct"
      },
      {
        "type": "field",
        "domain": ["geosite:category-ru", "geosite:apple"],
        "outboundTag": "direct"
      },
      {
        "type": "field",
        "ip": ["geoip:apple-ip"],
        "outboundTag": "direct"
      },
      {
        "type": "field",
        "ip": ["geoip:category-ru-ip", "geoip:category-by-ip", "geoip:category-kz-ip"],
        "outboundTag": "direct"
      }
    ]
  }
}
```

### Mihomo (Clash)

```yaml
rule-providers:
  ru-domain:
    type: http
    behavior: domain
    url: "https://github.com/this-xkit/flowvy_routing/releases/latest/download/category-ru.yaml"
    path: ./rules/ru-domain.yaml
    interval: 86400
  ru-ip:
    type: http
    behavior: ipcidr
    url: "https://github.com/this-xkit/flowvy_routing/releases/latest/download/category-ru-ip.yaml"
    path: ./rules/ru-ip.yaml
    interval: 86400
  apple:
    type: http
    behavior: domain
    url: "https://github.com/this-xkit/flowvy_routing/releases/latest/download/apple.yaml"
    path: ./rules/apple.yaml
    interval: 86400
  apple-ip:
    type: http
    behavior: ipcidr
    url: "https://github.com/this-xkit/flowvy_routing/releases/latest/download/apple-ip.yaml"
    path: ./rules/apple-ip.yaml
    interval: 86400

rules:
  - RULE-SET,ru-domain,DIRECT
  - RULE-SET,ru-ip,DIRECT
  - RULE-SET,apple,DIRECT
  - RULE-SET,apple-ip,DIRECT
```

## Источники данных

| Источник | Данные |
|----------|--------|
| [MetaCubeX/meta-rules-dat](https://github.com/MetaCubeX/meta-rules-dat) | category-ru, yandex, mailru, drweb, kaspersky, apple, private |
| [itdoginfo/allow-domains](https://github.com/itdoginfo/allow-domains) | Российские сервисы, доступные из-за рубежа |
| [hydraponique/roscomvpn-geosite](https://github.com/hydraponique/roscomvpn-geosite) | Коллекция российских доменов |
| [v2fly/geoip](https://github.com/v2fly/geoip) | IP-адреса по странам (RU, BY, KZ) |
| [ARIN/RIPE/APNIC](https://ipinfo.io/AS714) | IP-адреса Apple Inc. (AS714) |

## Особенности

- **Фильтрация по TLD** — из внешних источников берутся только домены в зонах `.ru`, `.su`, `.рф`, `.moscow`, `.tatar`
- **Дополнительные сервисы** — добавлены популярные российские сервисы в других зонах (yandex.com, vk.com и др.)
- **Keywords** — поддержка поиска по подстроке (avito, ozon, wildberries и др.)
- **Оптимизация для iOS** — компактные списки для работы в Network Extension (~50MB лимит)
- **Еженедельное обновление** — GitHub Actions обновляет списки каждый понедельник в 4:00 UTC

## Лицензия

MIT
