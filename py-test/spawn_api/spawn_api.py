import requests
import hmac
import hashlib


class Client:
    def __init__(self, client_id, secret):
        self.client_id = client_id
        self.secret = secret.encode('ascii')


class Device:
    def __init__(self, device_id, name):
        self.device_id = device_id
        self.name = name


class SpawnApi:
    def __init__(self, endpoint, client):
        self.client = client
        self.endpoint = endpoint

        self.auth_token = ""
        self.refresh_token = ""
        self.permissions = {}

    def signature(self, username, device_id):
        s = (self.client.client_id + device_id + username).encode('ascii')
        return hmac.new(self.client.secret, s, hashlib.sha512).hexdigest()

    def handle_login(self, resp):
        json = resp.json()
        if resp.status_code != 200:
            return json["error"]
        else:
            self.auth_token = json["auth_token"]
            if len(json["refresh_token"]) > 0:
                self.refresh_token = json["refresh_token"]
            self.permissions = json["permissions"]
            return None

    def sign_up(self, username, password, device, locale, lang):
        request = {
            "client_id": self.client.client_id,
            "device_id": device.device_id,
            "device_name": device.name,
            "username": username,
            "password": password,
            "locale": locale,
            "lang": lang,
            "signature": self.signature(username, device.device_id)
        }

        resp = requests.post(self.endpoint + '/auth/register', json=request)
        return self.handle_login(resp)

    def sign_in(self, username, password, device, locale, lang):
        request = {
            "client_id": self.client.client_id,
            "device_id": device.device_id,
            "device_name": device.name,
            "username": username,
            "password": password,
            "auth_type": "base",
            "locale": locale,
            "lang": lang,
            "signature": self.signature(username, device.device_id)
        }

        resp = requests.post(self.endpoint + '/auth/login', json=request)
        return self.handle_login(resp)

    def do_refresh_token(self):
        request = {
            "auth_token": self.auth_token,
            "refresh_token": self.refresh_token,
        }

        resp = requests.post(self.endpoint + '/auth/refresh_token', json=request)
        return self.handle_login(resp)

    def get_state(self):
        resp = requests.get(self.endpoint + "/user/state", headers={"Authorization": "Bearer " + self.auth_token})
        json = resp.json()

        if resp.status_code != 200:
            return True, json["error"]
        else:
            return False, json

    def logout(self):
        resp = requests.post(self.endpoint + "/user/logout", headers={"Authorization": "Bearer " + self.auth_token})
        json = resp.json()
        if resp.status_code != 200:
            return json["error"]
        else:
            return None

    def get_profile(self):
        resp = requests.get(self.endpoint + "/profile", headers={"Authorization": "Bearer " + self.auth_token})
        json = resp.json()
        if resp.status_code != 200:
            return True, json["error"]
        else:
            return False, json


# Client for test purposes
TEST_CLEINT = Client("client-test-01", "~_7|cjU^L?l5JI/jqN)S7|-I;=wz6<")
