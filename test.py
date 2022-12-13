import requests
from dotenv import dotenv_values

config = dotenv_values(".env")  # config = {"USER": "foo", "EMAIL": "foo@example.org"}


def loadSong():
    url = "https://api.spotify.com/v1/me/player/currently-playing"
    header = { "Authorization" : "Bearer " + config.get("TOKEN_user_read_recently_played") }
    r = requests.get(url, headers=header)
    data = r.json()
    # print(data["item"]["name"])
    # print(data["item"]["album"]["name"])
    # print(data["item"]["artists"][0]["name"])
    # print(str(data["progress_ms"])[:-3])
    # print(str(data["item"]["duration_ms"])[:-3])
    # print(data["is_playing"])
    return {
                "name": data["item"]["name"],
                "album": data["item"]["album"]["name"],
                "artist": data["item"]["artists"][0]["name"],
                "progress": str(data["progress_ms"])[:-3],
                "duration": str(data["item"]["duration_ms"])[:-3],
                "isPlaying":data["is_playing"],
                "trackURL": data["item"]["href"]
            }

print(loadSong())