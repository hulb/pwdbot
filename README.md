##### pwdbot
`pwdbot` is a telegram bot for manage my password.Thanks to [telebot](https://github.com/tucnak/telebot), it makes me very easy to create a telegram bot. `pwdbot` just saves account info(eg. username, password etc.) in plaintext, encryption is needed.The `pwdbot` is available at [pwdbot](https://t.me/passwdbot)(need telegram installed).

**NOTICE: The [pwdbot](https://t.me/passwdbot) now is deployed on my own server, therefor I can see the all content you input.**

##### Commands available
- `/new` `[length](optional)` return a random string in specified length
- `/save` `[password]` `[account name]` save password of the account
- `/update` `[acount name]` `.` `[property name]` `=` `[value]` update the specified property of the account name
- `/get` `[account name]` get detail of the account
- `/search` `[search key]` fuzzy search accounts that match the key
- `/rm` `[account name]` delete the account
- `/list` list all accounts

##### TODO
- [ ] account data encrypt
