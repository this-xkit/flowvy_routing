# Flowvy Routing

**Конфигурация Mihomo для Remnawave с умной маршрутизацией трафика.**

Заблокированные сервисы идут через прокси, всё остальное — напрямую. Без лишних настроек.

## Возможности

**Раздельная маршрутизация** — YouTube, Discord, мессенджеры и AI-сервисы вынесены в отдельные группы. Каждую можно переключить на другой прокси или отключить независимо.

**СНГ трафик напрямую** — домены и IP-адреса России, Беларуси и Казахстана идут в DIRECT. Включены Apple-сервисы и популярные российские приложения.

**Обход блокировок** — списки заблокированных доменов и IP из [re:filter](https://github.com/legiz-ru/mihomo-rule-sets) обновляются автоматически каждые 24 часа.

**CDN-подстраховка** — трафик к Cloudflare, Akamai, Amazon, Fastly и другим CDN направляется в группу «Прочие сервисы» (по умолчанию DIRECT, можно переключить на прокси).

**Прямое подключение** — торренты, VPN-клиенты (Tailscale, WireGuard, NetBird), удалённый доступ (AnyDesk, RustDesk, TeamViewer) и игры всегда идут напрямую.

## Группы прокси

| Группа | По умолчанию | Назначение |
|---|---|---|
| **Заблокированные сервисы** | Случайный сервер | Основная группа для заблокированного трафика |
| **YouTube** | Заблокированные сервисы | YouTube и связанные домены |
| **Discord** | Заблокированные сервисы | Домены, голосовые каналы и процесс Discord |
| **Мессенджеры** | Заблокированные сервисы | Telegram, WhatsApp |
| **AI** | Заблокированные сервисы | OpenAI, Gemini и [другие](https://github.com/MetaCubeX/meta-rules-dat/blob/meta/geo/geosite/category-ai-!cn.list) |
| **СНГ сервисы** | Без Proxy | RU/BY/KZ домены и IP, Apple, российские приложения |
| **Прочие сервисы** | Без Proxy | CDN и всё, что не попало в другие группы |

## Использование

Скопируйте содержимое [`config.yaml`](config.yaml) и вставьте в конфигурацию Remnawave.

## Источники правил

- [MetaCubeX/meta-rules-dat](https://github.com/MetaCubeX/meta-rules-dat) — geosite/geoip наборы
- [legiz-ru/mihomo-rule-sets](https://github.com/legiz-ru/mihomo-rule-sets) — re:filter, торренты, Discord voice, RU-приложения
- [PentiumB/CDN-RuleSet](https://github.com/PentiumB/CDN-RuleSet) — IP-диапазоны CDN-провайдеров
- [this-xkit/flowvy_routing](https://github.com/this-xkit/flowvy_routing/releases) — собственные наборы правил (RU, BY, KZ, Apple, YouTube, private)
