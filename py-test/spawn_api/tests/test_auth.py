import unittest
import spawn_api as spawn
from .helper import SpawnConn


class TestAuth(unittest.TestCase):

    def setUp(self):
        self.cn = SpawnConn()

    # Register new user
    # Should: user is logged in after registration, device is confirmed, permission by default
    def testRegister(self):
        username = self.cn.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"
        err = self.cn.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        # permissions

        # device is confirmed after registration
        self.assertEqual(True, self.cn.api.permissions["is_device_confirmed"])

        self.assertEqual(False, self.cn.api.permissions["is_2fa_required"])
        self.assertEqual(False, self.cn.api.permissions["is_email_confirmed"])
        self.assertEqual(False, self.cn.api.permissions["is_locked"])
        self.assertEqual(0, self.cn.api.permissions["scopes"])

    # Register new user and try to register user with the same name
    # Should: server throws 'user-already-exist' error
    def testRegisterAlreadyExists(self):
        username = self.cn.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.cn.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        err = self.cn.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNotNone(err)
        self.assertEqual("user-already-exist", err["reason"])

    # Register new user and login with same device
    # Should: only one session for (user / device) pair is allowed, first session invalidated
    def testLogin(self):
        username = self.cn.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.cn.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        # only one login for (username / device) is allowed
        cn2 = SpawnConn()
        cn2.username = username
        cn2.password = password
        cn2.device = device

        err = cn2.sign_in()
        self.assertIsNone(err)

        # first session invalidated
        is_error, err = self.cn.api.get_state()
        self.assertTrue(is_error)
        self.assertEqual("session-not-found", err["reason"])

    # Register new user and login with new device
    # Should: login success, device isn't confirmed
    def testLoginWithNewDevice(self):
        username = self.cn.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.cn.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        # login with new device -> login ok, but device unconfirmed
        err = self.cn.api.sign_in(username, password, spawn.Device("test-device-1-new", "test-device-1-name-new"), "ru", "es")
        self.assertIsNone(err)

        self.assertEqual(False, self.cn.api.permissions["is_device_confirmed"])

    # Try to login with wrong credentials
    # Should: server rejects login (user-unknown)
    def testWrongLogin(self):
        username = self.cn.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        # wrong password
        err = self.cn.api.sign_in(username, password + "111", device, "ru", "es")
        self.assertIsNotNone(err)
        self.assertEqual("user-unknown", err["reason"])

        # wrong username
        err = self.cn.api.sign_in(username + "111", password, device, "ru", "es")
        self.assertIsNotNone(err)
        self.assertEqual("user-unknown", err["reason"])

    # Register new user and refresh token
    # Should: new token != old token, old token is invalid
    def testRefreshToken(self):
        username = self.cn.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.cn.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        old_auth = self.cn.api.auth_token[:]

        err = self.cn.api.do_refresh_token()
        self.assertIsNone(err)

        # new token must be differ
        self.assertNotEqual(old_auth, self.cn.api.auth_token)

        # token is invalidated after refresh
        self.cn.api.auth_token = old_auth
        err = self.cn.api.do_refresh_token()
        self.assertIsNotNone(err)
        self.assertEqual("token-expired", err["reason"])


if __name__ == '__main__':
    unittest.main()