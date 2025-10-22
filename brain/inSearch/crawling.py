import sys
import requests
from bs4 import BeautifulSoup
import threading
import json

valueGo = sys.argv[1]
print(f"value received from Go -> {valueGo}")
lenCrawle = 8

results = []
lock = threading.Lock()

class Crawler:
    def __init__(self, url):
        self.url = url
        try:
            self.response = requests.get(url=self.url, timeout=10)
            self.soup = BeautifulSoup(self.response.text, "html.parser")
        except Exception as e:
            print(f"Error fetching {url}: {e}")
            self.response = None
            self.soup = None

    def getText(self):
        return self.response.text if self.response else ""

    def getLinks(self):
        if not self.soup:
            return []
        links = [a["href"] for a in self.soup.find_all("a", href=True)]
        return links

    def getTitle(self):
        if not self.soup:
            return ""
        return self.soup.title.text if self.soup.title else ""

    def getSubTitle(self):
        if not self.soup:
            return []
        subTitle = [i.text for j in range(3, 7) for i in self.soup.find_all(f"h{j}")]
        return subTitle

    def getParagraph(self):
        if not self.soup:
            return []
        para = [i.text for i in self.soup.find_all("p")]
        return para

    def getData(self):
        return {
            "urls": self.url,
            "text": self.getText(),
            "link": self.getLinks(),
            "title": self.getTitle(),
            "subtitle": self.getSubTitle(),
            "paragraph": self.getParagraph()
        }


def getUrls(ls=None):
    if ls is None:
        ls = []
    
    try:
        with open("UrlDemo.txt", "r") as file:
            for line in file:
                line = line.strip()
                if line:
                    url = line.replace("value", valueGo)
                    ls.append(url)
    except FileNotFoundError:
        print("UrlDemo.txt not found")
    
    return ls


def crawl_and_store(url):
    crawler = Crawler(url)
    data = crawler.getData()
    
    with lock:
        results.append(data)


ls = getUrls()

threads = []
for url in ls:
    t = threading.Thread(target=crawl_and_store, args=(url,))
    threads.append(t)

for t in threads:
    t.start()

for t in threads:
    t.join()

with open("regodt.json", "w", encoding="utf-8") as file:
    json.dump(results, file, indent=4, ensure_ascii=False)

if len(results) > lenCrawle:
    exit(1)

print(f"Crawling terminé ! {len(results)} URLs traitées")