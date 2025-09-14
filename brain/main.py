import os.path
import time

from crawler import CrawlManager
import threading
import subprocess
import pandas as pd
import os

url = {
    # 🔍 Moteurs de recherche
    "search": [
        "https://www.google.com",
        "https://www.bing.com",
        "https://www.yahoo.com",
        "https://www.duckduckgo.com",
        "https://www.baidu.com",
        "https://www.qwant.com",
        "https://www.yandex.com"
    ],

    # 👥 Réseaux sociaux
    "social_networks": [
        "https://fr.wikipedia.org/wiki/YouTube",
        "https://fr.wikipedia.org/wiki/Facebook",
        "https://fr.wikipedia.org/wiki/Instagram",
        "https://fr.wikipedia.org/wiki/Discord",
        "https://fr.wikipedia.org/wiki/Telegram",
        "https://fr.wikipedia.org/wiki/Snapchat",
        "https://fr.wikipedia.org/wiki/Reddit",
        "https://fr.wikipedia.org/wiki/GitHub",
        "https://fr.wikipedia.org/wiki/Stack_Overflow",
        "https://en.wikipedia.org/wiki/Replit",
        # ajouts
        "https://www.twitter.com",
        "https://www.tiktok.com",
        "https://www.whatsapp.com",
        "https://www.linkedin.com",
        "https://www.pinterest.com"
    ],

    # 🎥 Vidéo & Streaming
    "video_streaming": [
        "https://www.youtube.com",
        "https://www.netflix.com",
        "https://www.twitch.tv",
        "https://www.dailymotion.com",
        "https://www.vimeo.com",
        "https://www.primevideo.com",
        "https://www.disneyplus.com",
        "https://www.hulu.com",
        "https://www.crunchyroll.com",
        "https://www.allocine.fr"
    ],

    # 🛒 E-commerce & Shopping
    "e_commerce": [
        "https://www.amazon.com",
        "https://www.amazon.fr",
        "https://www.ebay.com",
        "https://www.alibaba.com",
        "https://www.walmart.com",
        "https://www.etsy.com",
        "https://www.fnac.com",
        "https://www.cdiscount.com"
    ],

    # 📰 Médias & Actualités
    "news": [
        "https://www.bbc.com",
        "https://www.cnn.com",
        "https://www.lemonde.fr",
        "https://www.nytimes.com",
        "https://www.theguardian.com",
        "https://www.france24.com",
        "https://www.reuters.com",
        "https://www.liberation.fr",
        "https://www.20minutes.fr"
    ],

    # 📚 Culture & Connaissances
    "culture": [
        "https://fr.wikipedia.org/wiki/Cin%C3%A9ma",
        "https://fr.wikipedia.org/wiki/Musique",
        "https://fr.wikipedia.org/wiki/Litt%C3%A9rature",
        "https://fr.wikipedia.org/wiki/Peinture",
        "https://fr.wikipedia.org/wiki/Architecture",
        "https://fr.wikipedia.org/wiki/Jeu_vid%C3%A9o",
        # ajouts
        "https://www.wikipedia.org",
        "https://www.quora.com",
        "https://stackoverflow.com",
        "https://www.khanacademy.org",
        "https://www.ted.com"
    ],

    # 🎮 Jeux & Divertissement
    "gaming": [
        "https://www.roblox.com",
        "https://www.epicgames.com",
        "https://store.steampowered.com",
        "https://www.minecraft.net",
        "https://www.ign.com",
        "https://www.gamespot.com",
        "https://www.jeuxvideo.com"
    ],

    # 💼 Business & Finance
    "economy": [
        "https://fr.wikipedia.org/wiki/Bitcoin",
        "https://fr.wikipedia.org/wiki/Bourse",
        "https://fr.wikipedia.org/wiki/Inflation",
        "https://fr.wikipedia.org/wiki/Fonds_mon%C3%A9taire_international",
        # ajouts
        "https://www.bloomberg.com",
        "https://www.forbes.com",
        "https://www.investopedia.com",
        "https://www.boursorama.com"
    ],

    # 💻 Tech & Développeurs
    "technology": [
        "https://fr.wikipedia.org/wiki/Intelligence_artificielle",
        "https://fr.wikipedia.org/wiki/Cybers%C3%A9curit%C3%A9",
        "https://fr.wikipedia.org/wiki/Blockchain",
        "https://fr.wikipedia.org/wiki/Cloud_computing",
        "https://fr.wikipedia.org/wiki/5G",
        # ajouts
        "https://www.github.com",
        "https://www.gitlab.com",
        "https://www.techcrunch.com",
        "https://www.wired.com",
        "https://www.theverge.com"
    ],

    # 🎶 Musique
    "music": [
        "https://fr.wikipedia.org/wiki/Rap",
        "https://fr.wikipedia.org/wiki/Rock",
        "https://fr.wikipedia.org/wiki/Musique_pop",
        "https://fr.wikipedia.org/wiki/Musique_classique",
        "https://fr.wikipedia.org/wiki/Drake_(rappeur)",
        "https://fr.wikipedia.org/wiki/Eminem",
        # ajouts
        "https://www.spotify.com",
        "https://www.deezer.com",
        "https://www.soundcloud.com",
        "https://www.apple.com/apple-music",
        "https://www.tidal.com",
        "https://www.booska-p.com",
        "https://www.raprnb.com",
        "https://genius.com"
    ],

    # 🎬 Films & Animations
    "movies": [
        "https://fr.wikipedia.org/wiki/Inception",
        "https://fr.wikipedia.org/wiki/Matrix_(film)",
        "https://fr.wikipedia.org/wiki/Avengers_(film)",
        # ajouts animation
        "https://fr.wikipedia.org/wiki/Cars_(film)",
        "https://fr.wikipedia.org/wiki/Toy_Story",
        "https://fr.wikipedia.org/wiki/Frozen",
        "https://fr.wikipedia.org/wiki/Le_Roi_lion",
        "https://fr.wikipedia.org/wiki/Shrek"
    ],
    "anime": [
        "https://fr.wikipedia.org/wiki/One_Piece",
        "https://fr.wikipedia.org/wiki/Naruto",
        "https://fr.wikipedia.org/wiki/Dragon_Ball",
        "https://fr.wikipedia.org/wiki/Attack_on_Titan",
        "https://fr.wikipedia.org/wiki/Demon_Slayer",
        "https://fr.wikipedia.org/wiki/Bleach_(manga)",
        "https://fr.wikipedia.org/wiki/Hunter_%C3%97_Hunter",
        "https://fr.wikipedia.org/wiki/Death_Note",
        "https://fr.wikipedia.org/wiki/My_Hero_Academia",
        "https://fr.wikipedia.org/wiki/Jujutsu_Kaisen"
    ],

    # ⚽ Sport
    "sports": [
        "https://fr.wikipedia.org/wiki/Football",
        "https://fr.wikipedia.org/wiki/Basket-ball",
        "https://fr.wikipedia.org/wiki/Tennis",
        "https://fr.wikipedia.org/wiki/Boxe",
        "https://fr.wikipedia.org/wiki/Arts_martiaux_mixtes",
        "https://fr.wikipedia.org/wiki/Danse",
        # clubs
        "https://fr.wikipedia.org/wiki/FC_Barcelone",
        "https://fr.wikipedia.org/wiki/Real_Madrid_Club_de_F%C3%BAtbol",
        "https://fr.wikipedia.org/wiki/Paris_Saint-Germain_Football_Club",
        "https://fr.wikipedia.org/wiki/Manchester_United_Football_Club",
        "https://fr.wikipedia.org/wiki/Bayern_Munich_(football)",
        # joueurs
        "https://fr.wikipedia.org/wiki/Lionel_Messi",
        "https://fr.wikipedia.org/wiki/Cristiano_Ronaldo",
        "https://fr.wikipedia.org/wiki/Neymar",
        "https://fr.wikipedia.org/wiki/Kylian_Mbapp%C3%A9",
        "https://fr.wikipedia.org/wiki/Zlatan_Ibrahimovi%C4%87",
        # ajouts sport sites
        "https://www.espn.com",
        "https://www.fifa.com",
        "https://www.uefa.com",
        "https://www.lequipe.fr",
        "https://www.nba.com"
    ],

    # 🍴 Cuisine
    "food": [
        "https://fr.wikipedia.org/wiki/Pizza",
        "https://fr.wikipedia.org/wiki/Sushi",
        "https://fr.wikipedia.org/wiki/Kebab",
        "https://fr.wikipedia.org/wiki/Tacos",
        "https://fr.wikipedia.org/wiki/P%C3%A2tes_alimentaires",
        # ajouts
        "https://www.marmiton.org",
        "https://www.750g.com",
        "https://cuisine.journaldesfemmes.fr",
        "https://www.bonappetit.com"
    ],

    # 🌍 Pays, villes & nature
    "countries": [
        "https://fr.wikipedia.org/wiki/France",
        "https://fr.wikipedia.org/wiki/Arm%C3%A9nie",
        "https://fr.wikipedia.org/wiki/United_States",
        "https://fr.wikipedia.org/wiki/Japan",
        "https://fr.wikipedia.org/wiki/Brazil"
    ],
    "cities": [
        "https://fr.wikipedia.org/wiki/Paris",
        "https://fr.wikipedia.org/wiki/New_York",
        "https://fr.wikipedia.org/wiki/Tokyo",
        "https://fr.wikipedia.org/wiki/Moscou",
        "https://fr.wikipedia.org/wiki/Duba%C3%AF"
    ],
    "nature": [
        "https://fr.wikipedia.org/wiki/For%C3%AAt",
        "https://fr.wikipedia.org/wiki/Fleuve",
        "https://fr.wikipedia.org/wiki/Montagne",
        "https://fr.wikipedia.org/wiki/Oc%C3%A9an",
        "https://fr.wikipedia.org/wiki/D%C3%A9sert"
    ],

    # 🐾 Animaux
    "animals": [
        "https://fr.wikipedia.org/wiki/Chat",
        "https://fr.wikipedia.org/wiki/Chien",
        "https://fr.wikipedia.org/wiki/Lion",
        "https://fr.wikipedia.org/wiki/Dauphin",
        "https://fr.wikipedia.org/wiki/Aigle"
    ],

    # 🧑‍💻 Programmation & outils
    "programming": [
        "https://fr.wikipedia.org/wiki/Python_(langage)",
        "https://fr.wikipedia.org/wiki/C++",
        "https://fr.wikipedia.org/wiki/Java_(langage)",
        "https://fr.wikipedia.org/wiki/JavaScript",
        "https://fr.wikipedia.org/wiki/PHP"
    ],
    "programmingTools": [
        "https://fr.wikipedia.org/wiki/Git",
        "https://fr.wikipedia.org/wiki/Docker_(logiciel)",
        "https://fr.wikipedia.org/wiki/Kubernetes",
        "https://fr.wikipedia.org/wiki/Visual_Studio_Code",
        "https://fr.wikipedia.org/wiki/PyCharm"
    ],

    # 🌌 Astronomie & espace
    "astronomy": [
        "https://fr.wikipedia.org/wiki/Univers",
        "https://fr.wikipedia.org/wiki/Syst%C3%A8me_solaire",
        "https://fr.wikipedia.org/wiki/Trou_noir",
        "https://fr.wikipedia.org/wiki/Voie_lact%C3%A9e",
        "https://fr.wikipedia.org/wiki/Mars_(plan%C3%A8te)"
    ],
    "spaceMissions": [
        "https://fr.wikipedia.org/wiki/Apollo_11",
        "https://fr.wikipedia.org/wiki/Station_spatiale_internationale",
        "https://fr.wikipedia.org/wiki/SpaceX",
        "https://fr.wikipedia.org/wiki/Programme_Artemis"
    ],

    # 📖 Histoire & philosophie
    "history": [
        "https://fr.wikipedia.org/wiki/World_War_I",
        "https://fr.wikipedia.org/wiki/World_War_II",
        "https://fr.wikipedia.org/wiki/Ancient_Rome",
        "https://fr.wikipedia.org/wiki/Ancient_Egypt",
        "https://fr.wikipedia.org/wiki/Cold_War",
        "https://fr.wikipedia.org/wiki/Histoire_de_l%27Arm%C3%A9nie",
        "https://fr.wikipedia.org/wiki/Histoire_de_la_France"
    ],
    "philosophy": [
        "https://fr.wikipedia.org/wiki/Platon",
        "https://fr.wikipedia.org/wiki/Aristote",
        "https://fr.wikipedia.org/wiki/Friedrich_Nietzsche",
        "https://fr.wikipedia.org/wiki/Emmanuel_Kant"
    ],

    # 📐 Maths & sciences
    "maths": [
        "https://fr.wikipedia.org/wiki/Alg%C3%A8bre",
        "https://fr.wikipedia.org/wiki/G%C3%A9om%C3%A9trie",
        "https://fr.wikipedia.org/wiki/Probabilit%C3%A9",
        "https://fr.wikipedia.org/wiki/Statistiques"
    ],
    "anatomy": [
        "https://fr.wikipedia.org/wiki/Corps_humain",
        "https://fr.wikipedia.org/wiki/Cerveau",
        "https://fr.wikipedia.org/wiki/C%C5%93ur",
        "https://fr.wikipedia.org/wiki/Poumon",
        "https://fr.wikipedia.org/wiki/Squelette"
    ],

    # ⚙️ Hardware / Software
    "hardware": [
        "https://fr.wikipedia.org/wiki/Nvidia",
        "https://fr.wikipedia.org/wiki/Advanced_Micro_Devices",
        "https://fr.wikipedia.org/wiki/Intel",
        "https://fr.wikipedia.org/wiki/Apple_(entreprise)",
        "https://fr.wikipedia.org/wiki/Qualcomm",
        "https://fr.wikipedia.org/wiki/Random-access_memory",
        "https://fr.wikipedia.org/wiki/Solid-state_drive",
        "https://fr.wikipedia.org/wiki/Hard_disk_drive",
        "https://fr.wikipedia.org/wiki/Motherboard",
        "https://fr.wikipedia.org/wiki/Computer_fan",
        "https://fr.wikipedia.org/wiki/Computer_case"
    ],
    "software": [
        "https://fr.wikipedia.org/wiki/Logiciel",
        "https://fr.wikipedia.org/wiki/Microsoft_Windows",
        "https://fr.wikipedia.org/wiki/Linux",
        "https://fr.wikipedia.org/wiki/MacOS",
        "https://fr.wikipedia.org/wiki/Android",
        "https://fr.wikipedia.org/wiki/IOS"
    ],

    # 🚘 Véhicules & marques
    "car": [
        "https://fr.wikipedia.org/wiki/BMW",
        "https://fr.wikipedia.org/wiki/Mercedes-Benz",
        "https://fr.wikipedia.org/wiki/Audi",
        "https://fr.wikipedia.org/wiki/Ferrari",
        "https://fr.wikipedia.org/wiki/Tesla_Motors"
    ],
    "vehicles": [
        "https://fr.wikipedia.org/wiki/Motorcycle",
        "https://fr.wikipedia.org/wiki/Car",
        "https://fr.wikipedia.org/wiki/Airplane",
        "https://fr.wikipedia.org/wiki/Train",
        "https://fr.wikipedia.org/wiki/Submarine"
    ],

    # 🧑‍⚖️ Politique & religion
    "politics": [
        "https://fr.wikipedia.org/wiki/Organisation_des_Nations_unies",
        "https://fr.wikipedia.org/wiki/Union_europ%C3%A9enne",
        "https://fr.wikipedia.org/wiki/D%C3%A9mocratie",
        "https://fr.wikipedia.org/wiki/Pr%C3%A9sident"
    ],
    "religion": [
        "https://fr.wikipedia.org/wiki/Christianisme",
        "https://fr.wikipedia.org/wiki/Islam",
        "https://fr.wikipedia.org/wiki/Juda%C3%AFsme",
        "https://fr.wikipedia.org/wiki/Bouddhisme",
        "https://fr.wikipedia.org/wiki/Hindouisme",
        "https://fr.wikipedia.org/wiki/Sikhisme",
        "https://fr.wikipedia.org/wiki/Tao%C3%AFsme",
        "https://fr.wikipedia.org/wiki/Confucianisme",
        "https://fr.wikipedia.org/wiki/Shinto",
        "https://fr.wikipedia.org/wiki/Zoroastrisme"
    ],

    # 📖 Mythologie
    "mythology": [
        "https://fr.wikipedia.org/wiki/Mythologie_grecque",
        "https://fr.wikipedia.org/wiki/Mythologie_%C3%A9gyptienne",
        "https://fr.wikipedia.org/wiki/Mythologie_nordique"
    ]
}




levier = True

LIMIT = 200

def crawleCategorie (categoris,urls):
    global levier

    for i in urls:
        FileLine = FileRaw("data.csv")
        if FileLine >= LIMIT:
            levier = False
            print (f"Limit de {LIMIT} a ete deppaser -> crawling arreter")
            return

        manager = CrawlManager([i],categoris,profondeurMax=10,limit = 2000)
        manager.start()



# def FileRaw (filePath):
#     if not os.path.exists(filePath):
#         return 0
#     with open (filePath , "r" , encoding = "utf-8") as file:
#         return sum(1 for _ in file)


def FileRaw(filePath):
    if not os.path.exists(filePath):
        return 0

    try:
        with open(filePath, "r", encoding="utf-8") as f:
            return sum(1 for _ in f)
    except UnicodeDecodeError:
        try:
            with open(filePath, "r", encoding="latin-1") as f:
                return sum(1 for _ in f)
        except UnicodeDecodeError:
            try:
                df = pd.read_csv(filePath, encoding="latin-1", on_bad_lines="skip")
                df.to_csv(filePath, index=False, encoding="utf-8-sig", quoting=1)  # quoting=csv.QUOTE_ALL
                print(f"⚡ CSV nettoyé et réécrit en UTF-8 : {filePath}")
                return len(df)
            except Exception as e:
                print(f"Impossible de nettoyer le CSV : {e}")
                return 0



def RunCrawler ():
    global levier

    fileLine = FileRaw("data.csv")

    if fileLine >= LIMIT:
        print("csv plein pas de crawling")
        levier = False
        return

    thread = []

    for key, value in url.items():
        if not levier:
            break

        t = threading.Thread(target=crawleCategorie, args=(key, value),)
        thread.append(t)
        t.start()

    for t in thread:
        t.join()


def RunGO ():
    global levier

    if not levier:
        try:
            subprocess.run(["go","run","main.go"],check = True)
        except subprocess.SubprocessError as E:
            print(f"Error run go from python -> {E}")



def ContinueCrawling():
    global levier
    with open("binary.bin", "wb") as file:
        file.write(b"\x01")
    levier = True




if __name__ == "__main__":
    while True:
        RunCrawler()
        RunGO()
        ContinueCrawling()