import unittest
from .helper import SpawnConn


class TestProfile(unittest.TestCase):

    def setUp(self):
        self.cn = SpawnConn()

    # Register new user; query user profile
    # Should: user profile is empty
    def testProfileForNewUser(self):
        err = self.cn.sign_up()
        self.assertIsNone(err)

        is_error, profile = self.cn.api.get_profile()
        self.assertFalse(is_error)

        self.assertEqual(self.cn.username, profile["auth_info"]["username"])

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
        err = self.cn.sign_up()
        self.assertIsNone(err)

        error = self.cn.api.update_personal_info("vasya", "pupkin", "1978-12-21")
        self.assertIsNone(error)

        error = self.cn.api.update_country("US")
        self.assertIsNone(error)

        is_error, profile = self.cn.api.get_profile()
        self.assertFalse(is_error)

        personal_info = profile["personal_info"]

        self.assertEqual('US', personal_info["country"])

        self.assertEqual('vasya', personal_info["first_name"])
        self.assertEqual('pupkin', personal_info["last_name"])
        self.assertEqual('1978-12-21', personal_info["birth_date"])


if __name__ == '__main__':
    unittest.main()
