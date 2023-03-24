package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

const baseURL = "https://www.okx.com/api/v5/market/index-tickers?instId="

type responseIndexTickers struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstId  string `json:"instId"`
		IdxPx   string `json:"idxPx"`
		High24h string `json:"high24h"`
		SodUtc0 string `json:"sodUtc0"`
		Open24h string `json:"open24h"`
		Low24h  string `json:"low24h"`
		SodUtc8 string `json:"sodUtc8"`
		Ts      string `json:"ts"`
	} `json:"data"`
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "token-pairs",
				DefaultText: "BTC-USDT",
			},
			&cli.StringFlag{
				Name:  "format",
				Value: "txt",
			},
		},

		Action: func(ctx *cli.Context) error {
			data := getData(ctx)
			if len(data) == 0 {
				return errors.New("无法获取数据")
			}

			// json to struct
			rITList := make([]responseIndexTickers, 0)
			for _, v := range data {
				r := responseIndexTickers{}
				err := json.Unmarshal(v, &r)
				if err != nil {
					log.Println("err:", err)
					continue
				}
				rITList = append(rITList, r)
			}

			// output json
			if ctx.String("format") == "json" {
				for _, v := range rITList {
					bList, _ := json.MarshalIndent(v, "", "    ")
					fmt.Println(string(bList))
				}

				return nil
			}

			// output txt
			msgList := make([]string, 0)
			for _, v := range rITList {
				msg := parseMsg(v)
				msgList = append(msgList, msg)
			}
			fmt.Println(strings.Join(msgList, "\n"))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getData(ctx *cli.Context) [][]byte {
	tokenPairs := strings.Split(ctx.String("token-pairs"), ",")
	data := make([][]byte, 0)

	for _, v := range tokenPairs {
		respBody := getResponse(v)
		data = append(data, respBody)
	}
	return data
}

func getResponse(tokenPairName string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+tokenPairName, nil)
	req.Header.Add("Context-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("err:", err)
		return []byte{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err:", err)
		return []byte{}
	}
	return body
}

func parseMsg(r responseIndexTickers) string {
	d := r.Data[0]

	example := "%s:\nIdxPx:%s\to24:%s\th24:%s\tl24:%s"
	return fmt.Sprintf(example, d.InstId, d.IdxPx, d.Open24h, d.High24h, d.Low24h)
}
