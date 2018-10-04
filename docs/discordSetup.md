How to execute:
``` Windows: ./go-supreme.exe testFileV2.json settings.json````
The files can be anywhere but the above example the two files are in the same folder as the exe. The setting file is opitonal and if not specified will use settings of (300,800,800).

Download: https://drive.google.com/open?id=1IdH-MMU7Bj3cWBqIvDkeA2_KtJDHCatT

The Taskfile is valid json (you can test that with this site https://jsonlint.com/):
```[
  {
    "taskName": "task1 - doesn't matter what you put here",
    "item": {
      "keywords": [
        "Briefs",
        "Boxer"
      ],
      "size": "medium",
      "color": "white",
      "category": "accessories"
    },
    "account": {
      "person": {
        "firstname": "Jax",
        "lastname": "Blax",
        "email": "none+0RU3@gmail.com",
        "phoneNumber": "354-143-9568"
      },
      "address": {
        "address1": "0RU3 123 HoneySuckle Ave",
        "address2": "",
        "zipcode": "85542",
        "city": "Springfield",
        "state": "WA",
        "country": "USA"
      },
      "card": {
        "cardtype": "notneeded",
        "number": "1234 2541 2154 5487",
        "month": "09",
        "year": "2022",
        "cvv": "789"
      }
    }
  }
]```
Currently the card type and taskName field is not used in processing. The taskName field exists so you can add a note to the task.

Settings file, also just json:
```{
  "refreshWait": 300,
  "atcWait": 800,
  "checkoutWait": 800
}```
All are in milliseconds, refreshWait is how long to wait before refreshing a category while monitoring or waiting for a restock.
atcWait is waiting before adding the item to the cart.
checkoutWait time to wait before completing the checkout.