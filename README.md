# Mihomo RU Config

🎯 **Конфигурация Mihomo, адаптированная под Рунет и популярные зарубежные сервисы.**

![FLClash 2](https://github.com/user-attachments/assets/97e14b01-9c53-467c-b24d-a18f3190266d)

Эта конфигурация оптимизирована для удобного, гибкого и безопасного доступа к российским сайтам, блокированным ресурсам, AI-сервисам, Discord, YouTube и т.д., с максимально прозрачной маршрутизацией.

## 💡 Особенности

- 📦 **Группы прокси по категориям**:
  - 🏠 RU сайты — весь российский и белорусский трафик направляется напрямую (DIRECT)
  - 🥷🏻 Глобальный Proxy — для заблокированных ресурсов
  - 🤖 AI — прокси для OpenAI, Gemini и других нейросетей
  - ▶️ YouTube & 💬 Discord — отдельная прокси-группа
  - 🔓 Без Proxy — ручной выбор прямого подключения

- 🌐 **Гибкая маршрутизация**:
  - Гео-наборы MetaCubeX
  - Правила для Cloudflare, торрентов, Discord, AI, RU/BY и прочих
  - Собственные inline-правила для RU-контента и блокируемых доменов

- 🛡️ **Безопасность и фильтрация**:
  - Встроенный adblock через `oisd_big`
  - Блокировка UDP DNS
  - Поддержка Fake-IP, sniffer, TUN и расширенного DNS

- ⚡ **Автовыбор быстрого прокси** (`url-test`)

## 🚀 Начало работы

Скопируйте конфиг ниже и вставьте его в вашу панель. Конфиг создавался под Remnawave.

<details>
<summary>📄 Показать/скрыть полный конфигурационный файл</summary>
  
```yaml
mixed-port: 7890
allow-lan: true
tcp-concurrent: true
enable-process: true
find-process-mode: always
mode: rule
log-level: info
ipv6: false
bind-address: "*"
keep-alive-interval: 30
unified-delay: false
profile:
  store-selected: true
  store-fake-ip: true
sniffer:
  enable: true
  force-dns-mapping: true
  parse-pure-ip: true
  sniff:
    HTTP:
      ports:
        - 80
        - 8080-8880
      override-destination: true
    TLS:
      ports:
        - 443
        - 8443
tun:
  enable: true
  stack: gvisor
  auto-route: true
  auto-detect-interface: true
  dns-hijack:
    - any:53
  strict-route: true
dns:
  enable: true
  prefer-h3: false
  use-hosts: true
  use-system-hosts: true
  listen: 127.0.0.1:6868
  ipv6: false
  enhanced-mode: redir-host
  default-nameserver:
    - tls://1.1.1.1
    - tls://1.0.0.1
    - tls://8.8.8.8
    - tls://8.8.4.4
    - tls://9.9.9.9
    - https://8.8.8.8/dns-query
    - https://1.1.1.1/dns-query
    - https://9.9.9.9/dns-query
  proxy-server-nameserver:
    - tls://1.1.1.1
    - tls://1.0.0.1
    - tls://8.8.8.8
    - tls://8.8.4.4
    - tls://9.9.9.9
    - https://8.8.8.8/dns-query
    - https://cloudflare-dns.com/dns-query
    - https://9.9.9.9/dns-query
  direct-nameserver:
    - tls://77.88.8.8#DIRECT
    - https://77.88.8.8/dns-query#DIRECT
  nameserver:
    - https://cloudflare-dns.com/dns-query#PROXY

proxies:

proxy-groups:
  - name: 🥷🏻 Глобальный Proxy (без RU/Torrent) # Используется для проксирования заблокированных в РФ сайтов, по умолчанию подключается к серверу с наименьшим пингом
    icon: https://cdn.jsdelivr.net/gh/Koolson/Qure@master/IconSet/Dark/Speedtest.png
    type: select
    proxies:
      - ⚡️ Самый быстрый
      - 🔓 Без Proxy
  - name: ▶️ YouTube & 💬 Discord # Используется для проксирования только YouTube и Discord,  по умолчанию подключается к серверу с наименьшим пингом
    icon: https://cdn.jsdelivr.net/gh/Koolson/Qure@master/IconSet/Dark/Game.png
    type: select
    proxies:
      - ⚡️ Самый быстрый
      - 🔓 Без Proxy
  - name: 🤖 AI # Используется для проксирования только AI-сервисов,  по умолчанию подключается к серверу с наименьшим пингом
    icon: https://cdn.jsdelivr.net/gh/Koolson/Qure@master/IconSet/Dark/Bot.png
    type: select
    proxies:
      - ⚡️ Самый быстрый
      - 🔓 Без Proxy
  - name: 🏠 RU сайты # Используется для проксирования RU и BY,  по умолчанию направляется в DIRECT
    icon: https://cdn.jsdelivr.net/gh/Koolson/Qure@master/IconSet/Color/Russia.png
    type: select
    proxies:
      - 🔓 Без Proxy
      - ⚡️ Самый быстрый
  - name: 🏳️ Other sites # Используется для проксирования трафика не попавшее ни в одно другое правило, по умолчанию использует селектор 🏠 RU сайты
    remnawave: # Кастомное поле используемое только в Remnawave (отключает добавление всех прокси в данную секцию, кроме указанных вручную)
      include-proxies: false
    hidden: true # селектор скрыт из интерфейса юзера
    type: select
    proxies:
      - 🏠 RU сайты
  - name: PROXY
    remnawave:
      include-proxies: false
    type: select
    hidden: true
    proxies:
      - 🥷🏻 Глобальный Proxy (без RU/Torrent)
  - name: 🔓 Без Proxy
    remnawave:
      include-proxies: false
    type: select
    hidden: true
    proxies:
      - DIRECT
  - name: ⚡️ Самый быстрый # Пингует хосты каждые 300 секунд и выбирает хост с наименьшим пингом
    hidden: true
    type: url-test
    tolerance: 150
    url: https://www.gstatic.com/generate_204
    interval: 300
    proxies:
  - name: ☁️ Cloudflare # Проксирует IP адреса Cloudflare, по умолчанию использует селектор ⚡️ Самый быстрый и скрыт от юзера, то есть Cloudflare проксируется всегда
    remnawave:
      include-proxies: false
    type: select
    hidden: true
    proxies:
      - ⚡️ Самый быстрый

rule-providers:
  ru-inline-banned:
    type: inline
    payload:
      - DOMAIN-SUFFIX,habr.com
      - DOMAIN-SUFFIX,kemono.su
      - DOMAIN-SUFFIX,jut.su
      - DOMAIN-SUFFIX,kara.su
      - DOMAIN-SUFFIX,theins.ru
      - DOMAIN-SUFFIX,tvrain.ru
      - DOMAIN-SUFFIX,echo.msk.ru
      - DOMAIN-SUFFIX,the-village.ru
      - DOMAIN-SUFFIX,snob.ru
      - DOMAIN-SUFFIX,novayagazeta.ru
      - DOMAIN-SUFFIX,moscowtimes.ru
      - DOMAIN-KEYWORD,animego
      - DOMAIN-KEYWORD,yummyanime
      - DOMAIN-KEYWORD,yummy-anime
      - DOMAIN-KEYWORD,animeportal
      - DOMAIN-KEYWORD,anime-portal
      - DOMAIN-KEYWORD,animedub
      - DOMAIN-KEYWORD,anidub
      - DOMAIN-KEYWORD,animelib
      - DOMAIN-KEYWORD,ikianime
      - DOMAIN-KEYWORD,anilibria
    behavior: classical
  ru-inline:
    type: inline
    payload:
      - DOMAIN-SUFFIX,2ip.ru
      - DOMAIN-SUFFIX,yastatic.net
      - DOMAIN-SUFFIX,yandex.net
      - DOMAIN-SUFFIX,yandex.kz
      - DOMAIN-SUFFIX,yandex.com
      - DOMAIN-SUFFIX,mycdn.me
      - DOMAIN-SUFFIX,vk.com
      - DOMAIN-SUFFIX,.ru
      - DOMAIN-SUFFIX,.su
      - DOMAIN-SUFFIX,.by
      - DOMAIN-SUFFIX,.ru.com
      - DOMAIN-SUFFIX,.ru.net
      - DOMAIN-SUFFIX,kudago.com
      - DOMAIN-SUFFIX,kinescope.io
      - DOMAIN-SUFFIX,redheadsound.studio
      - DOMAIN-SUFFIX,plplayer.online
      - DOMAIN-SUFFIX,lomont.site
      - DOMAIN-SUFFIX,remanga.org
      - DOMAIN-SUFFIX,shopstory.live
      - DOMAIN-KEYWORD,miradres
      - DOMAIN-KEYWORD,premier
      - DOMAIN-KEYWORD,shutterstock
      - DOMAIN-KEYWORD,2gis
      - DOMAIN-KEYWORD,diginetica
      - DOMAIN-KEYWORD,kinescopecdn
      - DOMAIN-KEYWORD,researchgate
      - DOMAIN-KEYWORD,springer
      - DOMAIN-KEYWORD,nextcloud
      - DOMAIN-KEYWORD,kaspersky
      - DOMAIN-KEYWORD,stepik
      - DOMAIN-KEYWORD,likee
      - DOMAIN-KEYWORD,snapchat
      - DOMAIN-KEYWORD,yappy
      - DOMAIN-KEYWORD,pikabu
      - DOMAIN-KEYWORD,okko
      - DOMAIN-KEYWORD,wink
      - DOMAIN-KEYWORD,kion
      - DOMAIN-KEYWORD,roblox
      - DOMAIN-KEYWORD,ozon
      - DOMAIN-KEYWORD,wildberries
      - DOMAIN-KEYWORD,aliexpress
    behavior: classical
  geosite-ru:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geosite/category-ru.mrs
    path: ./provider/rule-set/geosite-ru.mrs
    interval: 86400
  drweb:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geosite/drweb.mrs
    path: ./provider/rule-set/drweb.mrs
    interval: 86400
  geoip-ru:
    type: http
    behavior: ipcidr
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geoip/ru.mrs
    path: ./provider/rule-set/geoip-ru.mrs
    interval: 86400
  geoip-by:
    type: http
    behavior: ipcidr
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geoip/by.mrs
    path: ./provider/rule-set/geoip-by.mrs
    interval: 86400
  geosite-private:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geosite/private.mrs
    path: ./provider/rule-set/geosite-private.mrs
    interval: 86400
  geoip-private:
    type: http
    behavior: ipcidr
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geoip/private.mrs
    path: ./provider/rule-set/geoip-private.mrs
    interval: 86400
  discord_domains:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geosite/discord.mrs
    path: ./provider/rule-set/discord_domains.mrs
    interval: 86400
  discord_vc:
    type: inline
    payload:
      - AND,((IP-CIDR,138.128.136.0/21),(NETWORK,udp),(DST-PORT,50000-51000))
      - AND,((IP-CIDR,162.158.0.0/15),(NETWORK,udp),(DST-PORT,50000-51000))
      - AND,((IP-CIDR,172.64.0.0/13),(NETWORK,udp),(DST-PORT,50000-51000))
      - AND,((IP-CIDR,34.0.0.0/15),(NETWORK,udp),(DST-PORT,50000-51000))
      - AND,((IP-CIDR,34.2.0.0/15),(NETWORK,udp),(DST-PORT,50000-51000))
      - AND,((IP-CIDR,35.192.0.0/12),(NETWORK,udp),(DST-PORT,50000-51000))
      - AND,((IP-CIDR,35.208.0.0/12),(NETWORK,udp),(DST-PORT,50000-51000))
      - AND,((IP-CIDR,5.200.14.128/25),(NETWORK,udp),(DST-PORT,50000-51000))
      - AND,((IP-CIDR,66.22.192.0/18),(NETWORK,udp),(DST-PORT,50000-51000))
    behavior: classical
  refilter_domains:
    type: inline
    behavior: domain
    format: mrs
    url: https://github.com/legiz-ru/mihomo-rule-sets/raw/main/re-filter/domain-rule.mrs
    path: ./provider/rule-set/domain-rule.mrs
    interval: 86400
  refilter_ipsum:
    type: http
    behavior: ipcidr
    format: mrs
    url: https://github.com/legiz-ru/mihomo-rule-sets/raw/main/re-filter/ip-rule.mrs
    path: ./provider/rule-set/ip-rule.mrs
    interval: 86400
  youtube:
    type: http
    behavior: domain
    format: mrs
    interval: 86400
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geosite/youtube.mrs
    path: ./provider/rule-set/youtube.mrs
  oisd_big:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/legiz-ru/mihomo-rule-sets/raw/main/oisd/big.mrs
    path: ./provider/rule-set/oisd/big.mrs
    interval: 86400
  torrent-trackers:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/legiz-ru/mihomo-rule-sets/raw/main/other/torrent-trackers.mrs
    path: ./provider/rule-set/torrent-trackers.mrs
    interval: 86400
  torrent-clients:
    type: http
    behavior: classical
    format: yaml
    url: https://github.com/legiz-ru/mihomo-rule-sets/raw/main/other/torrent-clients.yaml
    path: ./provider/rule-set/torrent-clients.yaml
    interval: 86400
  ru-bundle:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/legiz-ru/mihomo-rule-sets/raw/main/ru-bundle/rule.mrs
    path: ./provider/rule-set/ru-bundle/rule.mrs
    interval: 86400
  full-vpn:
    type: inline
    behavior: classical
    payload:
      - NETWORK,tcp
      - NETWORK,udp
  openai:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geosite/openai.mrs
    path: ./provider/rule-set/openai.mrs
    interval: 86400
  gemini:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geosite/google-gemini.mrs
    path: ./provider/rule-set/gemini.mrs
    interval: 86400
  ai-bundle:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geosite/category-ai-!cn.mrs
    path: ./provider/rule-set/geo/geosite/category-ai-!cn.mrs
    interval: 86400
  cloudflare:
    type: http
    behavior: ipcidr
    format: mrs
    url: https://github.com/MetaCubeX/meta-rules-dat/raw/meta/geo/geoip/cloudflare.mrs
    path: ./provider/rule-set/cloudflare.mrs
    interval: 86400

rules:
  - RULE-SET,geosite-private,DIRECT,no-resolve
  - RULE-SET,geoip-private,DIRECT
  - AND,((NETWORK,udp),(DST-PORT,443)),REJECT
  - AND,((NETWORK,udp),(DST-PORT,853)),REJECT
  - RULE-SET,oisd_big,REJECT
  - OR,((RULE-SET,torrent-clients),(RULE-SET,torrent-trackers)),DIRECT
  - OR,((DOMAIN,ipwho.is),(DOMAIN,api.ip.sb),(DOMAIN,ipapi.co),(DOMAIN,ipinfo.io),(DOMAIN-SUFFIX,2ip.io),(DOMAIN-SUFFIX,2ipcore.com)),🥷🏻 Глобальный Proxy (без RU/Torrent)
  - OR,((RULE-SET,discord_domains),(RULE-SET,discord_vc),(PROCESS-NAME,Discord.exe),(RULE-SET,youtube)),▶️ YouTube & 💬 Discord
  - RULE-SET,ai-bundle,🤖 AI
  - RULE-SET,cloudflare,☁️ Cloudflare
  - RULE-SET,ru-bundle,🥷🏻 Глобальный Proxy (без RU/Torrent)
  - RULE-SET,refilter_domains,🥷🏻 Глобальный Proxy (без RU/Torrent)
  - RULE-SET,refilter_ipsum,🥷🏻 Глобальный Proxy (без RU/Torrent)
  - RULE-SET,ru-inline-banned,🥷🏻 Глобальный Proxy (без RU/Torrent)
  - RULE-SET,ru-inline,🏠 RU сайты
  - RULE-SET,geosite-ru,🏠 RU сайты
  - RULE-SET,geoip-ru,🏠 RU сайты
  - RULE-SET,geoip-by,🏠 RU сайты
  - RULE-SET,drweb,🏠 RU сайты
  - RULE-SET,full-vpn,🏳️ Other sites
  - MATCH,🥷🏻 Глобальный Proxy (без RU/Torrent)
```

</details>

---

> 🤍 Разработано с любовью к свободе интернета и логике маршрутизации.
