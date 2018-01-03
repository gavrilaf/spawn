import unittest
import spawn_api as spawn
import uuid


class TestAccounts(unittest.TestCase):

    def setUp(self):
        self.client = spawn.TEST_CLEINT
        self.endpoint = 'http://localhost:8080'
        self.api = spawn.SpawnApi(self.endpoint, self.client)

    @staticmethod
    def get_name():
        return str(uuid.uuid4()) + "@spawn.com"

    def testNewUser(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        is_error, accounts = self.api.get_accounts()
        self.assertFalse(is_error)

        # new user has empty accounts list
        self.assertEqual(0, len(accounts["accounts"]))


if __name__ == '__main__':
    unittest.main()
