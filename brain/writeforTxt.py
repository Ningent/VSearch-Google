import os
from dotenv import load_dotenv
import psycopg2

# Charger le .env
load_dotenv()

# Infos BDD
user = os.getenv("SupabaseName")
password = os.getenv("SupabaseMDP")
dbName = os.getenv("SupabaseDB")
host = os.getenv("SupabaseURL")
port = os.getenv("SupabasePort")
table = os.getenv("TableName")  # normalement "main"

# Connexion à la BDD
conn = psycopg2.connect(
    host=host,
    port=port,
    user=user,
    password=password,
    dbname=dbName,
    sslmode="require"
)

cur = conn.cursor()


cur.execute(f'SELECT "URL" FROM {table};')
urls = [row[0] for row in cur.fetchall()]
try:
    with open("LinkTxt.txt", "r", encoding="utf-8") as f:
        existing_links = set(f.read().splitlines())
except FileNotFoundError:
    existing_links = set()

# Nouveau contenu à ajouter
new_links = [
    "https://fr.wikipedia.org/wiki/Nvidia",
    "https://fr.wikipedia.org/wiki/Advanced_Micro_Devices",
    # ... autres liens
]

# Combiner et éviter les doublons
all_links = existing_links.union(new_links)

# Réécrire dans le fichier
with open("LinkTxt.txt", "w", encoding="utf-8") as f:
    for link in sorted(all_links):
        f.write(link + "\n")

