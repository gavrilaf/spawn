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

    # Register new user and query session state
    # Should return session state
    def testGetState(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        is_error, state = self.api.get_state()
        self.assertFalse(is_error, "Error is {}".format(state))

        self.assertEqual("es", state["lang"])
        self.assertEqual("ru", state["locale"])
        self.assertEqual(self.api.permissions, state["permissions"])

    # Register new user; logout; query session state
    # Should returns error (session-not-found)
    def testLogout(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        is_error, state = self.api.get_state()
        self.assertFalse(is_error, "Error is {}".format(state))

        err = self.api.logout()
        self.assertIsNone(err)

        # should be error
        is_error, err = self.api.get_state()
        self.assertTrue(is_error)
        self.assertEqual(err["reason"], "session-not-found")

    # Register new user; login with new device; query devices list
    # Should two devices in list; first is confirmed; second is current (active session)
    def testGetDevices(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "it", "es")
        self.assertIsNone(err)

        err = self.api.sign_in(username, password, spawn.Device("test-device-1-new", "test-device-1-name-new"), "ru",
                               "en")
        self.assertIsNone(err)

        is_error, devices = self.api.get_devices()
        self.assertFalse(is_error)

        devices = devices["devices"]
        self.assertEqual(2, len(devices))

        first = [x for x in devices if x["device_id"] == "test-device-1"]
        self.assertEqual(1, len(first))

        self.assertEqual("test-device-1-name", first[0]["device_name"])
        self.assertEqual(True, first[0]["is_confirmed"])
        self.assertEqual(False, first[0]["is_current"])
        self.assertEqual("it", first[0]["locale"])
        self.assertEqual("es", first[0]["lang"])
        self.assertIsNotNone(first[0]["login_ip"])
        self.assertIsNotNone(first[0]["login_region"])
        self.assertIsNotNone(first[0]["login_time"])
        self.assertIsNotNone(first[0]["user_agent"])

        current = [x for x in devices if x["device_id"] == "test-device-1-new"]
        self.assertEqual(1, len(current))

        self.assertEqual("test-device-1-name-new", current[0]["device_name"])
        self.assertEqual(False, current[0]["is_confirmed"])
        self.assertEqual(True, current[0]["is_current"])
        self.assertEqual("ru", current[0]["locale"])
        self.assertEqual("en", current[0]["lang"])
        self.assertIsNotNone(current[0]["login_ip"])
        self.assertIsNotNone(current[0]["login_region"])
        self.assertIsNotNone(current[0]["login_time"])
        self.assertIsNotNone(current[0]["user_agent"])

    # Register new user; login with new device; try to delete current device; try to delete first device
    # Should could not delete current device (delete-current-device),
    # after deleting first device only one device in devices list
    def testDeleteDevice(self):
        username = self.get_name()
        device = spawn.Device("device-1", "device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "en", "en")
        self.assertIsNone(err)

        err = self.api.sign_in(username, password, spawn.Device("device-2", "device-2-name"), "en", "en")
        self.assertIsNone(err)

        is_error, devices = self.api.get_devices()
        self.assertFalse(is_error)

        devices = devices["devices"]
        self.assertEqual(2, len(devices))

        err = self.api.delete_device("device-2")
        self.assertIsNotNone(err)
        self.assertEqual(err["reason"], "delete-current-device")

        err = self.api.delete_device("device-1")
        self.assertIsNone(err)

        is_error, devices = self.api.get_devices()
        self.assertFalse(is_error)

        devices = devices["devices"]
        self.assertEqual(1, len(devices))
        self.assertEqual("device-2", devices[0]["device_id"])
        self.assertEqual("device-2-name", devices[0]["device_name"])

    # Register new user; login with new device; delete first device; check first session state
    # Should device session is invalidated when device is deleted
    def testDeleteDeviceInvalidateSession(self):
        username = self.get_name()
        device = spawn.Device("device-1", "device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "en", "en")
        self.assertIsNone(err)

        # login with other device
        api2 = spawn.SpawnApi(self.endpoint, self.client)
        err = api2.sign_in(username, password, spawn.Device("device-2", "device-2-name"), "en", "en")
        self.assertIsNone(err)

        # session 2 is in valid state
        is_error, _ = api2.get_state()
        self.assertFalse(is_error)

        err = self.api.delete_device("device-2")
        self.assertIsNone(err)

        # session 2 is invalidated
        is_error, _ = api2.get_state()
        self.assertTrue(is_error)


if __name__ == '__main__':
    unittest.main()
