import spawn_api as spawn
import uuid


class SpawnConn:
    def __init__(self):
        self.client = spawn.TEST_CLEINT
        self.endpoint = 'http://localhost:8080'
        self.api = spawn.SpawnApi(self.endpoint, self.client)

        self.username = self.get_name()
        self.device = spawn.Device("test-device-1", "test-device-1-name")
        self.password = "password"
        self.locale = "ru"
        self.lang = "es"

    @staticmethod
    def get_name():
        return str(uuid.uuid4()) + "@spawn.com"

    def sign_up(self):
        return self.api.sign_up(self.username, self.password, self.device, self.locale, self.lang)

    def sign_in(self):
        return self.api.sign_in(self.username, self.password, self.device, self.locale, self.lang)