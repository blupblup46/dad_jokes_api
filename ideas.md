# Server

## global

1) au demarrage de serveur
- Appel API en mode batch (configurable: premier appel ON/OFF)(API batch: https://icanhazdadjoke.com/search)
- Génération de fichier JSON avec les données de API
- Récup des blagues perso
- fusion blagues perso/API 

# Client

- /getJoke                -> une blague au pif
- /getJoke/{id}           -> blague à l'id
- /getJokes               -> toutes les blagues
- /getJokes/from_api      -> toutes les blagues du fichier API
- /getJokes/from_perso    -> toutes les blagues du fichier API
- /initJokes              -> fait appel API et reset fichier JSON de blagues + régénére map de blagues fusionnées