**Download**
```
Mac: https://drive.google.com/open?id=1YlltpQDTpu7595FQjeCMBhpRFXF8Mv4b
Windows: https://drive.google.com/open?id=193IdAYbDmLxo9uaxbaKT8I8VuynXUM86
```
**How to Run**
In some terminal / command line navigate to the folder with the application is, your settings file and task file should be in the same folder if you want to use the below command.
```Mac: ./go-supreme testFileV2.json settings.json
Windows: ./go-supreme.exe testFileV2.json settings.json````
The files can be anywhere but the above example the two files are in the same folder as the app. The setting file is optional and if not specified will use settings of (300,800,800). I don't currently have a good set of recommended settings.
** Task File **
The Taskfile is valid json (you can test that with this site https://jsonlint.com/). An example file is provided, it is in the zip and called testFileV2.json
Currently the card type and taskName field is not used in processing. The taskName field exists so you can add a note to the task. The credit cards must be spaced correctly (Visa / Mastercards: XXXX XXXX XXXX XXXX Amex: XXXX XXXXXX XXXXX). The telephone number must have a dashes between them (584-530-4127).
Categories: "jackets", "shirts",  "tops/sweaters", "sweatshirts","pants", "t-shirts", "hats", "bags", "shorts", "accessories","skate", "shoes", "new" *NEW CATEGORY CAN ONLY BE USED BY MOBILE API*
Keywords: are Not CaSe SenSitive
API: "mobile" or "desktop"
** Settings **
Settings file is also json:
```{
  "refreshWait": 300,
  "atcWait": 800,
  "checkoutWait": 800
}```
All times are in milliseconds.
refreshWait is how long to wait before refreshing a category while monitoring or waiting for a restock.
atcWait is waiting before adding the item to the cart.
checkoutWait time to wait before completing the checkout.

**Warnings**
* I would not run for restocks with the mobile API