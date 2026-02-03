import requests
import yaml
import json
import os

# --- КОНФИГУРАЦИЯ ---
# 1. База v2fly (СНГ сегмент)
V2FLY_BASE_URL = "https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/"
V2FLY_CATEGORIES = ["category-ru", "category-ua", "category-by", "category-kz", "category-gov-ru"]

# 2. Внешние списки (сюда можно добавлять raw ссылки на чужие списки)
EXTERNAL_SOURCES = [
    # "https://raw.githubusercontent.com/username/repo/main/banks.txt"
]

# 3. Папки
LOCAL_FILES_DIR = "src"
DIST_DIR = "dist"
OUTPUT_NAME = "category-cis"

if not os.path.exists(DIST_DIR):
    os.makedirs(DIST_DIR)

def fetch_v2fly(files, processed=None):
    """Рекурсивно качает категории v2fly"""
    if processed is None: processed = set()
    domains = set()
    
    for filename in files:
        if filename in processed: continue
        processed.add(filename)
        
        try:
            resp = requests.get(f"{V2FLY_BASE_URL}{filename}")
            if resp.status_code != 200: continue
            
            for line in resp.text.splitlines():
                line = line.split('#')[0].strip()
                if not line: continue
                
                if line.startswith("include:"):
                    domains.update(fetch_v2fly([line.split(":")[1]], processed))
                elif not (line.startswith("regexp:") or line.startswith("keyword:")):
                    # Чистим от full:, domain: и оставляем чистый домен
                    d = line.replace("full:", "").replace("domain:", "").split("@")[0]
                    domains.add(d)
        except:
            print(f"Ошибка с {filename}")
    return domains

def fetch_external(urls):
    domains = set()
    for url in urls:
        try:
            resp = requests.get(url)
            if resp.status_code == 200:
                domains.update(l.strip() for l in resp.text.splitlines() if l.strip() and not l.startswith("#"))
        except:
            print(f"Ошибка с {url}")
    return domains

def read_local():
    domains = set()
    if not os.path.exists(LOCAL_FILES_DIR): return domains
    for f in os.listdir(LOCAL_FILES_DIR):
        if f.endswith(".txt"):
            with open(os.path.join(LOCAL_FILES_DIR, f), "r") as fl:
                domains.update(l.strip() for l in fl if l.strip() and not l.startswith("#"))
    return domains

def main():
    print("Сборка category-cis...")
    final_set = set()
    
    final_set.update(fetch_v2fly(V2FLY_CATEGORIES))
    final_set.update(fetch_external(EXTERNAL_SOURCES))
    final_set.update(read_local())
    
    sorted_list = sorted(list(final_set))
    
    # Mihomo
    with open(f"{DIST_DIR}/{OUTPUT_NAME}.yaml", "w") as f:
        yaml.dump({"payload": [f"DOMAIN-SUFFIX,{d}" for d in sorted_list]}, f, sort_keys=False)
        
    # Sing-box
    with open(f"{DIST_DIR}/{OUTPUT_NAME}.json", "w") as f:
        json.dump({"version": 1, "rules": [{"domain_suffix": sorted_list}]}, f, indent=2)

    print(f"Готово! Собрано {len(sorted_list)} доменов.")

if __name__ == "__main__":
    main()
