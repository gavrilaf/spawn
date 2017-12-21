import unittest
import spawn_api as spawn
import uuid


class TestProfile(unittest.TestCase):

    def setUp(self):
        self.client = spawn.TEST_CLEINT
        self.endpoint = 'http://localhost:8080'
        self.api = spawn.SpawnApi(self.endpoint, self.client)

    @staticmethod
    def get_name():
        return str(uuid.uuid4()) + "@spawn.com"

    def testRegister(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"
        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        # permissions

        # device is confirmed after registration
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
        self.assertEqual(err["reason"], "user-already-exist")

    def testLogin(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        # only one login for (username / device) is allowed
        err = self.api.sign_in(username, password, device, "ru", "es")
        self.assertIsNotNone(err)
        self.assertEqual(err["reason"], "session-already-exist")

        # logout
        err = self.api.logout()
        self.assertIsNone(err)

        # now you can sign in
        err = self.api.sign_in(username, password, device, "ru", "es")
        self.assertIsNone(err)

        # device is confirmed
        self.assertEqual(self.api.permissions["is_device_confirmed"], True)

    def testLoginWithNewDevice(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        # login with new device -> login ok, but device unconfirmed
        err = self.api.sign_in(username, password, spawn.Device("test-device-1-new", "test-device-1-name-new"), "ru", "es")
        self.assertIsNone(err)

        self.assertEqual(self.api.permissions["is_device_confirmed"], False)

    def testWrongLogin(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        # wrong password
        err = self.api.sign_in(username, password + "111", device, "ru", "es")
        self.assertIsNotNone(err)
        self.assertEqual(err["reason"], "user-unknown")

        # wrong username
        err = self.api.sign_in(username + "111", password, device, "ru", "es")
        self.assertIsNotNone(err)
        self.assertEqual(err["reason"], "user-unknown")

    def testRefreshToken(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        old_auth = self.api.auth_token[:]
        old_refresh = self.api.refresh_token[:]

        err = self.api.do_refresh_token()
        self.assertIsNone(err)

        #self.assertFalse(self.api.auth_token == old_auth) # TODO: Fix after token nonce task



if __name__ == '__main__':
    unittest.main()