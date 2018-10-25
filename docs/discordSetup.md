**Download**
```
Mac: https://mega.nz/#!fWBnQKrK!V8d6KrZ18BQ-sPhz0Mktns4yd82UMg5wGJBocejzU2Y
Windows: https://mega.nz/#!HDBDlKia!s7ETygH8Hxf1WhyhV-Dlqdkzb6oRcY9VA_nohhb7n6Y
```
**How to Run**
In some terminal / command line navigate to the folder with the application is, your settings file and task file should be in the same folder if you want to use the below command.
```Mac: ./supreme testFileV3.json settings.json
Windows: ./supreme.exe testFileV3.json settings.json````
The files can be anywhere but the above example the two files are in the same folder as the app. The setting file is optional and if not specified will use settings of (300,800,800). I don't currently have a good set of recommended settings.
** Task File **
The Taskfile is valid json (you can test that with this site https://jsonlint.com/). An example file is provided, it is in the zip and called testFileV3.json
Currently the card type and taskName field is not used in processing. The taskName field exists so you can add a note to the task. The credit cards must be spaced correctly (Visa / Mastercards: XXXX XXXX XXXX XXXX Amex: XXXX XXXXXX XXXXX). The telephone number must have a dashes between them (584-530-4127).
Categories: "jackets", "shirts",  "tops/sweaters", "sweatshirts","pants", "t-shirts", "hats", "bags", "shorts", "accessories","skate", "shoes"
Keywords: are Not CaSe SenSitive
API: "mobile" or "desktop"
** Settings **
Settings file is also json:
```{
  "startTime": "2018-10-18T14:59:30.000Z",
  "refreshWait": 300,
  "atcWait": 800,
  "checkoutWait": 800
}```
startTime will start all tasks after the specific size.
All times are in milliseconds.
refreshWait is how long to wait before refreshing a category while monitoring or waiting for a restock.
atcWait is waiting before adding the item to the cart.
checkoutWait time to wait before completing the checkout.
You can see the per task settings in the profiled example task file.