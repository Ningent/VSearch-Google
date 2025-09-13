import pandas as pd
import time

filePath = "data.xlsx"

empty_df = pd.DataFrame(columns=["URL", "title", "subTitle", "paragraphe", "lien", "categorie"])

empty_df.to_excel(filePath, index=False)

print(f"Le fichier {filePath} a été vidé avec succès.")
time.sleep(5000)
exit(1)
