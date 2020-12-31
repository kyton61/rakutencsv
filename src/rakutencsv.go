package main

import (
				"encoding/csv"
				"fmt"
				"io"
				"log"
				"os"
				"strconv"
				"strings"
				"time"

				iconv "github.com/djimenez/iconv-go"
)

// csv生データを格納する構造体
type TradeRaw struct {
				ContractDate	string
				ValueDate	string
				Ticker	string
				Name	string
				BuyOrSell	string
				Account	string
				NumOfShare	int
				UnitSellPrice float64
				SellPrice	float64
				BuyPriceAvg	float64
				ProfitAndLoss	float64
				ProfitRatio	float64
}

// 月別に加工したデータを格納する構造体
type TradeMonthData struct {
				YYYYmm	string
				ProfitAvg	float64
				LossAvg	float64
				ProfitMax	float64
				LossMax	float64
				ProfitRatioAvg	float64
				LossRatioAvg	float64
				ProfitRatioMax	float64
				LossRatioMax	float64
				WinTrade	int
				LossTrade	int
				WinRatio	float64
				TradeNum	int
}

// TradeRaw構造体のレシーバ：値のセットメソッド
func (t *TradeRaw) Set(contractDate string, valueDate string, ticker string, name string, buyOrSell string,
account string, numOfShare string, unitSellPrice string, sellPrice string, buyPriceAvg string, profitAndLoss string) {
				var i int
				var f float64
				var err error

				// debug
				fmt.Println(contractDate)
				fmt.Println(valueDate)
				fmt.Println(ticker)
				fmt.Println(name)
				fmt.Println(buyOrSell)
				fmt.Println(account)
				fmt.Println(numOfShare)
				fmt.Println(unitSellPrice)
				fmt.Println(sellPrice)
				fmt.Println(buyPriceAvg)
				fmt.Println(profitAndLoss)

				// 値のセット
				t.ContractDate = contractDate
				t.ValueDate = valueDate
				t.Ticker = ticker
				t.Name = name
				t.BuyOrSell = buyOrSell
				t.Account = account

				// string型以外のデータは型変換してセット
				i, err = strconv.Atoi(numOfShare)
				if err != nil {
								log.Fatal(err)
				}
				t.NumOfShare = i
				// 数字中のカンマを削除
				unitSellPrice = strings.Replace(unitSellPrice, ",", "", -1)
				f, err = strconv.ParseFloat(unitSellPrice, 32)
				if err != nil {
								log.Fatal(err)
				}
				t.UnitSellPrice = f
				sellPrice = strings.Replace(sellPrice, ",", "", -1)
				f, err = strconv.ParseFloat(sellPrice, 32)
				if err != nil {
								log.Fatal(err)
								}
				t.SellPrice = f
				buyPriceAvg = strings.Replace(buyPriceAvg, ",", "", -1)
				f, err = strconv.ParseFloat(buyPriceAvg, 32)
				if err != nil {
								log.Fatal(err)
				}
				t.BuyPriceAvg = f
				profitAndLoss = strings.Replace(profitAndLoss, ",", "", -1)
				f, err = strconv.ParseFloat(profitAndLoss, 32)
				if err != nil {
								log.Fatal(err)
				}
				t.ProfitAndLoss = f
				// 生データから利益率を計算
				t.ProfitRatio = t.ProfitAndLoss / t.SellPrice * 100

}

// "2006/01/02" 形式を "200601" 形式に変換
func ContractDateToYYYYmm(strdate string) string {
				var YYYYmm string
				t, _ := time.Parse("2006/01/02", strdate)
				YYYYmm = t.Format("200601")

				return YYYYmm
}

// TradeMonthData構造体のレシーバ：値のセットメソッド
func (t *TradeMonthData) Set(YYYYmm string, profitAndLoss float64, profitRatio float64) {
				// 初期化
				t.YYYYmm = YYYYmm
				t.ProfitAvg = 0
				t.LossAvg = 0
				t.ProfitMax = 0
				t.LossMax = 0
				t.ProfitRatioAvg = 0
				t.LossRatioAvg = 0
				t.ProfitRatioMax = 0
				t.LossRatioMax = 0
				t.WinTrade = 0
				t.LossTrade = 0
				t.WinRatio = 0
				t.TradeNum = 1

				// 利益率がプラスの場合
				if profitRatio >= 0 {
								t.ProfitAvg = profitAndLoss
								t.ProfitMax = profitAndLoss
								t.ProfitRatioAvg = profitRatio
								t.ProfitRatioMax = profitRatio
								t.WinTrade++
				// 利益率がマイナスの場合
				} else {
								t.LossAvg = profitAndLoss
								t.LossMax = profitAndLoss
								t.LossRatioAvg = profitRatio
								t.LossRatioMax = profitRatio
								t.LossTrade++
				}
}

// TradeMonthData構造体のレシーバ：値の追加メソッド
func (t *TradeMonthData) Add(tradeMonthDatas []TradeMonthData, tradeRaw TradeRaw) []TradeMonthData {
				YYYYmm := ContractDateToYYYYmm(tradeRaw.ContractDate)
				// 月ごとの統計データの中身と挿入するデータを比較する
				// 月データが配列内に存在する場合とそうでない場合でデータ挿入方法を変更する
				for i, tmd := range tradeMonthDatas {
								// 月データがすでに存在する場合
								if tmd.YYYYmm == YYYYmm {
												// 利益率のプラスマイナス判定
												if tradeRaw.ProfitRatio >= 0 {
																// 利益率がプラスの場合
																tradeMonthDatas[i].ProfitAvg += tradeRaw.ProfitAndLoss
																if tradeMonthDatas[i].ProfitMax < tradeRaw.ProfitAndLoss {
																				tradeMonthDatas[i].ProfitMax = tradeRaw.ProfitAndLoss
																}
																tradeMonthDatas[i].ProfitRatioAvg += tradeRaw.ProfitRatio
																if tradeMonthDatas[i].ProfitRatioMax < tradeRaw.ProfitRatio {
																				tradeMonthDatas[i].ProfitRatioMax = tradeRaw.ProfitRatio
																}
																tradeMonthDatas[i].WinTrade++
																tradeMonthDatas[i].TradeNum++
												} else {
																// 利益率がマイナスの場合
																tradeMonthDatas[i].LossAvg += tradeRaw.ProfitAndLoss
                                if tradeMonthDatas[i].LossMax > tradeRaw.ProfitAndLoss {
                                        tradeMonthDatas[i].LossMax = tradeRaw.ProfitAndLoss
                                }
                                tradeMonthDatas[i].LossRatioAvg += tradeRaw.ProfitRatio
                                if tradeMonthDatas[i].LossRatioMax < tradeRaw.ProfitRatio {
                                        tradeMonthDatas[i].LossRatioMax = tradeRaw.ProfitRatio
                                }
                                tradeMonthDatas[i].LossTrade++
                                tradeMonthDatas[i].TradeNum++
												}
												return tradeMonthDatas
								}
				}
				// 月データが存在しない場合
				var tmd TradeMonthData
				tmd.Set(YYYYmm, tradeRaw.ProfitAndLoss, tradeRaw.ProfitRatio)
				tradeMonthDatas = append(tradeMonthDatas, tmd)
				return tradeMonthDatas
}

func main() {

				// csvファイルのオープン
				// TODO:読み込みパスを指定できるようにする
				f, err := os.Open("file.csv")
				if err != nil {
								log.Fatal(err)
				}

				r := csv.NewReader(f)

				var trs []TradeRaw
				var tmd TradeMonthData
        var tmds []TradeMonthData

				// ヘッダの読み飛ばし
				r.Read()

				// csvファイルの読み込み
				for {
								record, err := r.Read()
								if err == io.EOF {
												break
								}
								if err != nil {
												log.Fatal(err)
								}
								fmt.Println(record)
								fmt.Println(len(record))
								// フッター読み飛ばしのため、空行がでてきたらbreak
								if len(record[0]) == 0 {
												break
								}

								// ヘッダを除いてcsvファイルデータを構造体に読み込む
								var tr TradeRaw
								tr.Set(record[0],record[1],record[2],record[3],record[4],record[5],record[6],record[7],record[8],record[9],record[10])
								trs = append(trs, tr)
				}

				// debug
				for _, v := range trs {
								println(v.Name, v.ProfitRatio)
				}

				// csvファイルのrawデータを月別のデータに格納
				for _, tr := range trs {
								tmds = tmd.Add(tmds, tr)
				}

				// 格納した月別データの平均と勝率を計算
				for i, v := range tmds {
								if v.WinTrade != 0 {
												tmds[i].ProfitAvg = v.ProfitAvg / float64(v.WinTrade)
												tmds[i].ProfitRatioAvg = v.ProfitRatioAvg / float64(v.WinTrade)
												tmds[i].WinRatio = float64(v.WinTrade) / float64(v.TradeNum) * 100
								} else if v.LossTrade != 0 {
												tmds[i].LossAvg = v.LossAvg / float64(v.LossTrade)
												tmds[i].LossRatioAvg = v.LossRatioAvg / float64(v.LossTrade)
								}
				}



				// debug
				println("==== 加工データの表示 =====")
				for _, v := range tmds {
								fmt.Printf("YYYYmm: %s\n", v.YYYYmm)
								fmt.Printf("平均利益[$]: %.2f\n", v.ProfitAvg)
								fmt.Printf("平均損失[$]: %.2f\n", v.LossAvg)
								fmt.Printf("最大利益[$]: %.2f\n", v.ProfitMax)
								fmt.Printf("最大損失[$]: %.2f\n", v.LossMax)
								fmt.Printf("平均利益率[%%]: %.2f\n", v.ProfitRatioAvg)
								fmt.Printf("平均損失率[%%]: %.2f\n", v.LossRatioAvg)
								fmt.Printf("最大利益率[%%]: %.2f\n", v.ProfitRatioMax)
								fmt.Printf("最大損失率[%%]: %.2f\n", v.LossRatioMax)
								fmt.Printf("勝ちトレード数: %d\n", v.WinTrade)
								fmt.Printf("負けトレード数: %d\n", v.LossTrade)
								fmt.Printf("勝率[%%]: %.2f\n", v.WinRatio)
								fmt.Printf("総トレード数: %d\n", v.TradeNum)
				}

				// O_WRONLY:書き込みモード開く, O_CREATE:無かったらファイルを作成
				fw, err := os.OpenFile("./output.csv", os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
								log.Fatal(err)
				}

				defer fw.Close()

				// ファイルを空にする
				err = fw.Truncate(0)

				converter, err := iconv.NewWriter(fw, "utf-8", "sjis")
				if err != nil {
								log.Fatal(err)
				}

				writer := csv.NewWriter(converter)

				// ヘッダの追記
				writer.Write([]string{"年月", "平均利益", "平均損失", "最大利益", "最大損失",
				"平均利益率", "平均損失率", "最大利益率", "最大損失率", "勝ちトレード数",
				"負けトレード数", "勝率", "総トレード数"})

				for _, v := range tmds {
								writer.Write([]string{v.YYYYmm,
								strconv.FormatFloat(v.ProfitAvg, 'f', 2, 64),
								strconv.FormatFloat(v.LossAvg, 'f', 2, 64),
								strconv.FormatFloat(v.ProfitMax, 'f', 2, 64),
								strconv.FormatFloat(v.LossMax, 'f', 2, 64),
								strconv.FormatFloat(v.ProfitRatioAvg, 'f', 2, 64),
								strconv.FormatFloat(v.LossRatioAvg, 'f', 2, 64),
								strconv.FormatFloat(v.ProfitRatioMax, 'f', 2, 64),
								strconv.FormatFloat(v.LossRatioMax, 'f', 2, 64),
								strconv.Itoa(v.WinTrade),
								strconv.Itoa(v.LossTrade),
								strconv.FormatFloat(v.WinRatio, 'f', 2, 64),
								strconv.Itoa(v.TradeNum)})
				}
				writer.Flush()
}

