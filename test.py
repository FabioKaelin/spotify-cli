import requests
from dotenv import dotenv_values
import base64

config = dotenv_values(".env")  # config = {"USER": "foo", "EMAIL": "foo@example.org"}


# def loadSong():
#     url = "https://accounts.spotify.com/authorize?"
#     url += "response_type=code&"
#     url += "client_id="+config.get("CLIENT_ID") +"&"
#     url += "redirect_uri=localhost:9876&"
#     url += "scope=user-read-currently-playing"
#     # header = { "Authorization" : "Bearer " + config.get("TOKEN_user_read_recently_played") }
#     r = requests.get(url)
#     # r = requests.get(url, headers=header)
#     data = r.json()
#     with open("b.json", "w") as f:
#         f.write(data)
# loadSong()


# var request = require('request'); // "Request" library

# var client_id = 'CLIENT_ID'; // Your client id
# var client_secret = 'CLIENT_SECRET'; // Your secret


# authOptions = {
#   "url": 'https://accounts.spotify.com/api/token',
#   "headers": {
#     'Authorization': 'Basic ' + str(base64.b64encode((config.get("CLIENT_ID")+":"+config.get("CLIENT_SECRET")).encode("utf-8")))
#   },
#   "form": {
#     "grant_type": 'client_credentials'
#   },
#   "json": True
# }

# str(
#     base64.b64encode(
#         (config.get("CLIENT_ID")+":"+config.get("CLIENT_SECRET")).encode("ascii")
#         )
#     )

# print(type(base64.b64encode(
#         (config.get("CLIENT_ID")+":"+config.get("CLIENT_SECRET")).encode("ascii")
#         )))

the_data = str({"grant_type": 'client_credentials'})
the_data = {"grant_type": 'client_credentials'}
headers = {'Content-Type': 'application/x-www-form-urlencoded','Authorization': 'Basic ' + (base64.b64encode((config.get("CLIENT_ID")+":"+config.get("CLIENT_SECRET")).encode("utf-8"))).decode("utf-8") }

# requests.post("http://bla.bla.example.com", data=the_data, headers=headers)
r = requests.post("https://accounts.spotify.com/api/token", data=the_data, headers=headers)
print(r.status_code)
# print(r.headers)
print(r.content)
# request.post(authOptions, function(error, response, body) {
#   if (!error && response.statusCode === 200) {

#     // use the access token to access the Spotify Web API
#     var token = body.access_token;
#     var options = {
#       url: 'https://api.spotify.com/v1/users/jmperezperez',
#       headers: {
#         'Authorization': 'Bearer ' + token
#       },
#       json: true
#     };
#     request.get(options, function(error, response, body) {
#       console.log(body);
#     });
#   }
# });