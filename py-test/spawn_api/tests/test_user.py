import unittest
import spawn_api as spawn
from .helper import SpawnConn


class TestUser(unittest.TestCase):

    def setUp(self):
        self.cn = SpawnConn()

    # Register new user and query session state
    # Should: return session state
    def testGetState(self):
        err = self.cn.sign_up()
        self.assertIsNone(err)

        is_error, state = self.cn.api.get_state()
        self.assertFalse(is_error, "Error is {}".format(state))

        self.assertEqual("es", state["lang"])
        self.assertEqual("ru", state["locale"])
        self.assertEqual(self.cn.api.permissions, state["permissions"])

    # Register new user; logout; query session state
    # Should: returns error (session-not-found)
    def testLogout(self):
        err = self.cn.sign_up()
        self.assertIsNone(err)

        err = self.cn.api.logout()
        self.assertIsNone(err)

        # should be error
        is_error, err = self.cn.api.get_state()
        self.assertTrue(is_error)
        self.assertEqual("session-not-found", err["reason"])

    # Register new user; login with new device; query devices list
    # Should: two devices in list; first is confirmed; second is current (active session)
    def testGetDevices(self):
        err = self.cn.sign_up()
        self.assertIsNone(err)

        cn2 = SpawnConn()
        cn2.username = self.cn.username
        cn2.password = self.cn.password
        cn2.device = spawn.Device("test-device-1-new", "test-device-1-name-new")
        cn2.locale = "us"
        cn2.lang = "it"

        err = cn2.sign_in()
        self.assertIsNone(err)

        is_error, devices = cn2.api.get_devices()
        self.assertFalse(is_error)

        devices = devices["devices"]
        self.assertEqual(2, len(devices))

        first = [x for x in devices if x["device_id"] == self.cn.device.device_id]
        self.assertEqual(1, len(first))

        self.assertEqual(self.cn.device.name, first[0]["device_name"])
        self.assertEqual(True, first[0]["is_confirmed"])
        self.assertEqual(False, first[0]["is_current"])
        self.assertEqual("ru", first[0]["locale"])
        self.assertEqual("es", first[0]["lang"])
        self.assertIsNotNone(first[0]["login_ip"])
        self.assertIsNotNone(first[0]["login_region"])
        self.assertIsNotNone(first[0]["login_time"])
        self.assertIsNotNone(first[0]["user_agent"])

        current = [x for x in devices if x["device_id"] == "test-device-1-new"]
        self.assertEqual(1, len(current))

        self.assertEqual("test-device-1-name-new", current[0]["device_name"])
        self.assertEqual(False, current[0]["is_confirmed"])  # new device, not confirmed
        self.assertEqual(True, current[0]["is_current"])  # active session
        self.assertEqual("us", current[0]["locale"])
        self.assertEqual("it", current[0]["lang"])
        self.assertIsNotNone(current[0]["login_ip"])
        self.assertIsNotNone(current[0]["login_region"])
        self.assertIsNotNone(current[0]["login_time"])
        self.assertIsNotNone(current[0]["user_agent"])

    # Register new user; login with new device; try to delete current device; try to delete first device
    # Should: could not delete current device (delete-current-device),
    # after deleting first device only one device in devices list
    def testDeleteDevice(self):
        err = self.cn.sign_up()
        self.assertIsNone(err)

        cn2 = SpawnConn()
        cn2.username = self.cn.username
        cn2.password = self.cn.password
        cn2.device = spawn.Device("device-2", "device-2-name")

        err = cn2.sign_in()
        self.assertIsNone(err)

        is_error, devices = self.cn.api.get_devices()
        self.assertFalse(is_error)

        devices = devices["devices"]
        self.assertEqual(2, len(devices))

        err = cn2.api.delete_device("device-2") # try to delete form device-2 session
        self.assertIsNotNone(err)
        self.assertEqual(err["reason"], "delete-current-device")

        err = self.cn.api.delete_device("device-2")  # try to delete form device-1 session
        self.assertIsNone(err)

        is_error, devices = self.cn.api.get_devices()
        self.assertFalse(is_error)

        devices = devices["devices"]
        self.assertEqual(1, len(devices))
        self.assertEqual(self.cn.device.device_id, devices[0]["device_id"])
        self.assertEqual(self.cn.device.name, devices[0]["device_name"])

    # Register new user; login with new device; delete first device; check first session state
    # Should: device session is invalidated when device is deleted
    def testDeleteDeviceInvalidateSession(self):
        err = self.cn.sign_up()
        self.assertIsNone(err)

        cn2 = SpawnConn()
        cn2.username = self.cn.username
        cn2.password = self.cn.password
        cn2.device = spawn.Device("device-2", "device-2-name")

        err = cn2.sign_in()
        self.assertIsNone(err)

        # session 2 is in valid state
        is_error, _ = cn2.api.get_state()
        self.assertFalse(is_error)

        err = self.cn.api.delete_device("device-2")
        self.assertIsNone(err)

        # session 2 is invalidated
        is_error, _ = cn2.api.get_state()
        self.assertTrue(is_error)


if __name__ == '__main__':
    unittest.main()
