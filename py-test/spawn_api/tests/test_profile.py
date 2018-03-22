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

    # Register new user; query user profile
    # Should user profile is empty
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
        self.assertEqual('', personal_info["first_name"])
        self.assertEqual('', personal_info["last_name"])
        self.assertEqual('', personal_info["country"])
        self.assertEqual('1800-01-01', personal_info["birth_date"]) # Empty date of birth constant

        phone_number = personal_info["phone_number"]
        self.assertEqual(False, phone_number["is_phone_confirmed"])
        self.assertEqual(0, phone_number["phone_country_code"])
        self.assertEqual('', phone_number["phone_number"])

    # Register new user; update profile & country
    def testUpdateUserPersonalInfo(self):
        username = self.get_name()
        device = spawn.Device("test-device-1", "test-device-1-name")
        password = "password"

        err = self.api.sign_up(username, password, device, "ru", "es")
        self.assertIsNone(err)

        error = self.api.update_personal_info("vasya", "pupkin", "1978-12-21")
        self.assertIsNone(error)

        error = self.api.update_country("US")
        self.assertIsNone(error)

        is_error, profile = self.api.get_profile()
        self.assertFalse(is_error)

        personal_info = profile["personal_info"]

        self.assertEqual('US', personal_info["country"])

        self.assertEqual('vasya', personal_info["first_name"])
        self.assertEqual('pupkin', personal_info["last_name"])
        self.assertEqual('1978-12-21', personal_info["birth_date"])



if __name__ == '__main__':
    unittest.main()
