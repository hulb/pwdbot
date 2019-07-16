##### pwdbot
pwdbot is a telegram bot for manage my password.Thanks to [telebot](https://github.com/tucnak/telebot), it make me very easy to make a telegram bot.pwdbot just save account info(eg. username, password etc.) in plaintext, encryption is needed.The pwdbot is available at [pwdbot](https://t.me/passwdbot)(need telegram installed).

**NOTICE: the [pwdbot](https://t.me/passwdbot) now is deplyed on my own server, therefor I can see the all content you input.**

##### commands available
- `/new` `length` return a random string in specified length
- `/save` `password` `account name` save password of the account
- `/update` `acount name` `property name``::``property value` update the specified property of the account name
- `/get` `account name` get detail of the account
- `/list` list all accounts

##### todo
- [ ] account data encrypt
- [ ] `/rm` delete the specified account
- [ ] `/search` fuzzy search 
