## how to use
### Download realized gains and losses csv from manex
- Go to マネックス証券(https://www.monex.co.jp/)
- Sign in to your account
- Go to page ログイン＞MY PAGE＞保有残高・口座管理＞売却損益明細
- download csv file from CSVダウンロード button

### Download exe file from git
- download manexcsv.exe
- if you build exe from src, please make cross compaile
	ex. env GOOS=windows GOARCH=amd64 go build manexcsv.go

### Calculate profits and losses
- put avobe files in same directory
- change csv file name to file.csv
- double click manexcsv.exe
- you can get output.csv in same directory
