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

    # Register new user; query accounts list
    # Accounts list is empty for new user
    def testProfileForNewUser(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        is_error, profile = self.api.get_profile()
        self.assertFalse(is_error)
        #print(profile)

        self.assertEqual(username, profile["auth_info"]["username"])

        # Empty personal info for new user
        personal_info = profile["personal_info"]
        self.assertEqual(personal_info["first_name"], '')
        self.assertEqual(personal_info["last_name"], '')
        self.assertEqual(personal_info["country"], '')
        self.assertEqual(personal_info["birth_date"], '1800-01-01T00:00:00Z') # Empty date of birth constant

        phone_number = personal_info["phone_number"]
        self.assertEqual(phone_number["is_phone_confirmed"], False)
        self.assertEqual(phone_number["phone_country_code"], 0)
        self.assertEqual(phone_number["phone_number"], '')


if __name__ == '__main__':
    unittest.main()
