import unittest
from .helper import SpawnConn

class TestAccounts(unittest.TestCase):

    def setUp(self):
        self.cn = SpawnConn()

    # Register new user; query accounts list
    # Accounts list is empty for new user
    def testNewUser(self):
        err = self.cn.sign_up()
        self.assertIsNone(err)

        is_error, accounts = self.cn.api.get_accounts()
        self.assertFalse(is_error)

        # new user has empty accounts list
        self.assertEqual(0, len(accounts["accounts"]))


if __name__ == '__main__':
    unittest.main()
