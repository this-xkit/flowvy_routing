# Flowvy Routing

Правила маршрутизации для Xray, V2Ray и Mihomo (Clash).

Автоматически собирает `geosite.dat`, `geoip.dat` и YAML файлы из нескольких источников.

## Скачать

### Xray / V2Ray
```
https://github.com/this-xkit/flowvy_routing/releases/latest/download/geosite.dat
https://github.com/this-xkit/flowvy_routing/releases/latest/download/geoip.dat
```

### Mihomo (Clash)
```
https://github.com/this-xkit/flowvy_routing/releases/latest/download/category-cis-domain.yaml
https://github.com/this-xkit/flowvy_routing/releases/latest/download/category-cis-ip.yaml
```

## Категории

| Категория | Описание |
|-----------|----------|
| `category-cis-domain` | Все домены СНГ (RU, BY, KZ + AntiFilter + Re-filter) |
| `category-cis-ip` | Все IP адреса СНГ (AntiFilter + Re-filter) |

## Примеры конфигурации

### Xray / V2Ray

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

### Mihomo (Clash)

```yaml
rule-providers:
  cis-domain:
    type: http
    behavior: domain
    url: "https://github.com/this-xkit/flowvy_routing/releases/latest/download/category-cis-domain.yaml"
    path: ./rules/cis-domain.yaml
    interval: 86400
  cis-ip:
    type: http
    behavior: ipcidr
    url: "https://github.com/this-xkit/flowvy_routing/releases/latest/download/category-cis-ip.yaml"
    path: ./rules/cis-ip.yaml
    interval: 86400

rules:
  - RULE-SET,cis-domain,DIRECT
  - RULE-SET,cis-ip,DIRECT
```

## Источники данных

- [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community) — category-ru, category-by, category-kz, category-gov-ru
- [v2fly/geoip](https://github.com/v2fly/geoip) — IP геолокация
- [1andrevich/Re-filter-lists](https://github.com/1andrevich/Re-filter-lists) — заблокированные домены/IP
- [community.antifilter.download](https://community.antifilter.download/) — сообщество AntiFilter

## Автоматическое обновление

GitHub Actions собирает и публикует релизы:
- Еженедельно (воскресенье, 4:00 UTC)
- При изменении файлов
- Вручную через GitHub Actions UI
