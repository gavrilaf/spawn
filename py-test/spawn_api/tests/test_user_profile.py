import unittest
import spawn_api as spawn
import uuid


class TestProfile(unittest.TestCase):

    def setUp(self):
        self.client = spawn.TEST_CLEINT
        self.api = spawn.SpawnApi('http://localhost:8080', self.client)

    def get_name(self):
        return str(uuid.uuid4()) + "@spawn.com"

    def testRegister(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"
        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        # permissions
        self.assertEqual(self.api.permissions["is_device_confirmed"], True)
        self.assertEqual(self.api.permissions["is_2fa_reqired"], False)
        self.assertEqual(self.api.permissions["is_email_confirmed"], False)
        self.assertEqual(self.api.permissions["is_locked"], False)
        self.assertEqual(self.api.permissions["scopes"], 0)

    def testRegisterAlreadyExists(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNotNone(err)
        self.assertEqual(err["scope"], "auth")
        self.assertEqual(err["reason"], "user-already-exist")

    def testLogin(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        err = self.api.sign_in(username, password, device, "ru", "es")
        self.assertIsNone(err)

        # device is confirmed
        self.assertEqual(self.api.permissions["is_device_confirmed"], True)

        # login with new device -> login ok, but device unconfirmed
        err = self.api.sign_in(username, password, spawn.Device("test-device-1-new", "test-device-1-name-new"), "ru", "es")
        self.assertIsNone(err)

        self.assertEqual(self.api.permissions["is_device_confirmed"], False)

        # wrong password
        err = self.api.sign_in(username, password + "111", device, "ru", "es")
        self.assertIsNotNone(err)
        self.assertEqual(err["scope"], "auth")
        self.assertEqual(err["reason"], "user-unknown")

        # wrong username
        err = self.api.sign_in(username + "111", password, device, "ru", "es")
        self.assertIsNotNone(err)
        self.assertEqual(err["scope"], "auth")
        self.assertEqual(err["reason"], "user-unknown")


if __name__ == '__main__':
    unittest.main()