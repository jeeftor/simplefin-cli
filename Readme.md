# Go CLI For SimpleFIN

simplefin.org is a website that will let you pull down finance data. This is a very simple Go client to pull your data down.

## Usage

At a minimum you must pass in your `access url` which is something you can generate through the simplefin/bridge interface.

```bash
sf --url https://3C9A2074C8EB44A694B197CA7ED9BDB15C3D2B89A74A4012F8AA22A9AF1DBDE4:EEF2D56BB0F1DB9E05A666E15213A7B2BD695E1F0D509D1D6708EAE39D6B6418@beta-bridge.simplefin.org/simplefin

```

or 

```bash
# Export Env Var
export SF_URL=https://3C9A2074C8EB44A694B197CA7ED9BDB15C3D2B89A74A4012F8AA22A9AF1DBDE4:EEF2D56BB0F1DB9E05A666E15213A7B2BD695E1F0D509D1D6708EAE39D6B6418@beta-bridge.simplefin.org/simplefin

# run the sf command
sf
```

## Output
And you'll see something like this

```bash
┌──────────────────────┬──────────────────────────────┬────────────┬──────────┬───┐
│ ORG NAME             │ ACCOUNT NAME                 │ BALANCE    │ CURRENCY │   │
├──────────────────────┼──────────────────────────────┼────────────┼──────────┼───┤
│ BANK1                │ Checking                     │ 13431.50   │ USD      │   │
│ BANK1                │ Visa® Card by VISA           │ -1786.24   │ USD      │   │
│ My Happy Investments │ Some Cash Account            │ 5264.92    │ USD      │   │
│ My Happy Investments │ A credit card.               │ -5404.82   │ USD      │   │
│ Crypto.              │ LOTS OF MONEY.               │ 9918791.79 │ USD      │ ⚠ │
└──────────────────────┴──────────────────────────────┴────────────┴──────────┴───┘
```

a ` ⚠` indiciates there may be an issue with the account