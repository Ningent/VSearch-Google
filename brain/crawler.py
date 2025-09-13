import requests
from bs4 import BeautifulSoup as bs4
import re
import json
import pandas as pd
import os
from collections import deque
import csv

# -------------------------------------------------
# Classe Crawler : s'occupe d'une seule page
# -------------------------------------------------
class Crawler:
    def __init__(self, url, categorie, visited_urls):
        self.url = url
        self.categorie = categorie
        self.visited_urls = visited_urls  # ensemble des URLs déjà crawlées

        self.soup = self.getSoup()
        self.titre = self.getTitle()
        self.subtitle = self.getSubtitles()
        self.paragraphe = self.getParagraphs()
        self.lien = self.getLinks()

        self.saveData()
        self.linkDataPrint()
        self.saveUrl()

        print (f" self.visitedUrls -> {self.visited_urls}")

    # 1. Récupération de la page HTML
    def getSoup(self):
        headers = {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
                          "AppleWebKit/537.36 (KHTML, like Gecko) "
                          "Chrome/91.0.4472.124 Safari/537.36"
        }
        try:
            response = requests.get(self.url, headers=headers )
            print("Statut HTTP :", response.status_code, "->", self.url)
            return bs4(response.text, "lxml")
        except Exception as e:
            print("Erreur requête:", e, "->", self.url)
            return None

    # 2. Extraction des infos
    def getTitle(self):
        return self.soup.title.text if self.soup and self.soup.title else None

    def getSubtitles(self):
        subtitles = []
        if not self.soup:
            return subtitles
        for tag in self.soup.find_all(["h2", "h3"]):
            span = tag.find("span", class_="mw-headline")
            if span:
                subtitles.append(span.text)
            else:
                subtitles.append(tag.text.strip())
        return subtitles

    def getParagraphs(self):
        return [p.text for p in self.soup.find_all("p")] if self.soup else []

    def getLinks(self):
        liens = []
        if not self.soup:
            return liens
        for a in self.soup.find_all("a"):
            i = a.get("href")
            if not i or not isinstance(i, str):
                continue

            if re.match(r"^https?://", i):
                liens.append(i)
            else:
                if i.startswith("/wiki") or i.startswith("/w"):
                    i = "https://wikipedia.org" + i
                elif i.startswith("#") or i == "#":
                    continue
                else:
                    i = "https:/" + i
                liens.append(i)
        return liens

    def linkDataPrint(self):
        print(f"URL {self.url} → {len(self.lien)} liens trouvés")

    # 3. Sauvegarde en CSV et Excel
    def saveData(self):
        csvPath = "data.csv"
        linkPath = "LinkTxt.txt"

        # --- Lire LinkTxt.txt pour savoir quelles URLs sont déjà dans le CSV ---
        visited_from_file = set()
        if os.path.exists(linkPath):
            with open(linkPath, "r", encoding="utf-8") as f:
                visited_from_file = set(line.strip() for line in f.readlines())

        # --- Préparer les données ---
        data = pd.DataFrame({
            "URL": [self.url],
            "title": [self.titre],
            "subTitle": [json.dumps(self.subtitle, ensure_ascii=False)],
            "paragraphe": [json.dumps(self.paragraphe, ensure_ascii=False)],
            "lien": [json.dumps(self.lien, ensure_ascii=False)],
            "categorie": [self.categorie]
        })

        # --- Ajouter la ligne au CSV uniquement si elle n’existe pas déjà ---
        if self.url not in visited_from_file:
            if os.path.exists(csvPath) and os.path.getsize(csvPath) > 0:
                try:
                    df_existing = pd.read_csv(csvPath, encoding="utf-8-sig", quoting=csv.QUOTE_ALL)
                except pd.errors.EmptyDataError:
                    df_existing = pd.DataFrame()
                df_combined = pd.concat([df_existing, data], ignore_index=True)
            else:
                df_combined = data


            # --- Écrire le CSV en échappant tout correctement ---
            df_combined.to_csv(csvPath, index=False, encoding="utf-8-sig", quoting=csv.QUOTE_ALL)

            print(f"Données ajoutées -> {self.url}")

            # --- Ajouter l’URL dans LinkTxt.txt pour ne pas la réécrire ---
            with open(linkPath, "a", encoding="utf-8") as f:
                f.write(self.url + "\n")
        else:
            print(f"{self.url} déjà présent dans LinkTxt.txt, skip CSV mais continue crawling")

    # 4. Sauvegarde des URLs crawlées
    def saveUrl(self):
        self.visited_urls.add(self.url)
        with open("LinkTxt.txt", "w", encoding="utf-8") as f:
            for u in self.visited_urls:
                f.write(u + "\n")




# -------------------------------------------------
# Classe CrawlManager : gère plusieurs pages
# -------------------------------------------------
class CrawlManager:
    def __init__(self, start_urls, categorie, profondeurMax, limit=200):
        self.queue = deque([(url, 0) for url in start_urls])
        self.visited = set()
        self.categorie = categorie
        self.profondeurMax = profondeurMax
        self.limit = limit

        if os.path.exists("LinkTxt.txt"):
            with open("LinkTxt.txt", "r", encoding="utf-8") as f:
                self.visited = set([line.strip() for line in f.readlines()])

    def start(self):
        while self.queue and self.limit > 0:
            url, depth = self.queue.popleft()
            if url in self.visited or depth > self.profondeurMax:
                continue
            try:
                crawler = Crawler(url, self.categorie, self.visited)
                self.limit -= 1
                for link in crawler.lien:
                    if link not in self.visited:
                        self.queue.append((link, depth + 1))
            except Exception as e:
                print("Erreur:", e, "->", url)

