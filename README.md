# Flowvy Routing

Правила маршрутизации для Xray, V2Ray и Mihomo (Clash).

Автоматически собирает `geosite.dat` и `geoip.dat` из нескольких источников.

## Использование

### Скачать последние файлы

```
# geosite.dat (домены)
https://github.com/this-xkit/flowvy_routing/releases/latest/download/geosite.dat

# geoip.dat (IP адреса)
https://github.com/this-xkit/flowvy_routing/releases/latest/download/geoip.dat
```

### Категории в geosite.dat

| Категория | Описание |
|-----------|----------|
| `category-cis-domain` | Все домены СНГ (RU, BY, KZ + AntiFilter + Re-filter) |
| `category-cis-ip` | Все IP адреса СНГ (AntiFilter + Re-filter) |

### Пример конфигурации Xray/V2Ray

```json
{
  "routing": {
    "domainStrategy": "IPIfNonMatch",
    "rules": [
      {
        "type": "field",
        "domain": ["geosite:category-cis-domain"],
        "outboundTag": "direct"
      },
      {
        "type": "field",
        "ip": ["geosite:category-cis-ip"],
        "outboundTag": "direct"
      }
    ]
  }
}
```

### Пример для Mihomo (Clash)

```yaml
rules:
  - GEOSITE,category-cis-domain,DIRECT
  - GEOIP,RU,DIRECT
```

## Источники данных

- [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community) — category-ru, category-by, category-kz, category-gov-ru
- [v2fly/geoip](https://github.com/v2fly/geoip) — IP геолокация
- [1andrevich/Re-filter-lists](https://github.com/1andrevich/Re-filter-lists) — заблокированные домены/IP
- [community.antifilter.download](https://community.antifilter.download/) — сообщество AntiFilter

## Автоматическое обновление

GitHub Actions собирает и публикует релизы:
- Еженедельно (воскресенье, 4:00 UTC)
- При изменении файлов в `data/`, `main.go`, или workflow
- Вручную через GitHub Actions UI
