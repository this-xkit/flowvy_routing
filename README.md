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

# geoip для заблокированных IP в России
https://github.com/this-xkit/flowvy_routing/releases/latest/download/geoip-refilter.dat
```

### Категории в geosite.dat

| Категория | Описание |
|-----------|----------|
| `category-ru` | Русские домены (прямой доступ) |
| `category-ua` | Украинские домены |
| `category-by` | Белорусские домены |
| `category-kz` | Казахстанские домены |
| `category-gov-ru` | Государственные сайты РФ |
| `category-blocked-ru` | Заблокированные в РФ домены (через прокси) |
| `private` | Приватные/локальные домены |
| `youtube`, `telegram`, `discord`, etc. | Популярные сервисы |

### Пример конфигурации Xray/V2Ray

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
        "domain": ["geosite:category-ru", "geosite:category-gov-ru"],
        "outboundTag": "direct"
      },
      {
        "type": "field",
        "domain": ["geosite:category-blocked-ru", "geosite:youtube"],
        "outboundTag": "proxy"
      }
    ]
  }
}
```

## Источники данных

- [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community) - категории доменов
- [v2fly/geoip](https://github.com/v2fly/geoip) - IP геолокация
- [1andrevich/Re-filter-lists](https://github.com/1andrevich/Re-filter-lists) - заблокированные домены/IP
- [community.antifilter.download](https://community.antifilter.download/) - сообщество AntiFilter

## Сборка локально

```bash
# Установить Go 1.22+
go mod download
go run main.go --datapath=./data --outputdir=./
```

## Автоматическое обновление

GitHub Actions собирает и публикует релизы:
- Еженедельно (воскресенье, 4:00 UTC)
- При изменении файлов в `data/`, `main.go`, или workflow
- Вручную через GitHub Actions UI
