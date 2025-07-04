package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// 股票数据结构
type StockData struct {
	Symbol    string  `json:"symbol"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
	Timestamp int64   `json:"timestamp"`
}

// API响应结构
type APIResponse struct {
	Data []StockData `json:"data"`
}

func main() {
	// 1. 配置股票代码和API密钥（示例使用Alpha Vantage API）
	symbol := "IBM" // 替换为你想分析的股票代码
	apiKey := "YOUR_API_KEY" // 替换为你的API密钥

	// 2. 获取今日股票数据
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", symbol, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("请求API失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		os.Exit(1)
	}

	// 3. 解析JSON数据
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("解析JSON失败: %v\n", err)
		os.Exit(1)
	}

	// 4. 提取今日数据（假设API返回最新数据是今天）
	today := time.Now().Format("2006-01-02")
	dailyData, ok := result["Time Series (Daily)"].(map[string]interface{})
	if !ok {
		fmt.Println("API返回数据格式不符合预期")
		os.Exit(1)
	}

	todayData, ok := dailyData[today].(map[string]interface{})
	if !ok {
		fmt.Println("未找到今日交易数据")
		os.Exit(1)
	}

	// 5. 分析数据
	open, _ := todayData["1. open"].(string)
	high, _ := todayData["2. high"].(string)
	low, _ := todayData["3. low"].(string)
	closePrice, _ := todayData["4. close"].(string)
	volume, _ := todayData["5. volume"].(string)

	// 转换为数值
	openVal := parseFloat(open)
	highVal := parseFloat(high)
	lowVal := parseFloat(low)
	closeVal := parseFloat(closePrice)
	volumeVal := parseInt(volume)

	// 6. 计算分析指标
	priceChange := closeVal - openVal
	percentChange := (priceChange / openVal) * 100
	volatility := (highVal - lowVal) / openVal * 100

	// 7. 输出分析结果
	fmt.Printf("\n=== %s 今日股票分析 (%s) ===\n", symbol, today)
	fmt.Printf("开盘价: %.2f\n", openVal)
	fmt.Printf("收盘价: %.2f\n", closeVal)
	fmt.Printf("最高价: %.2f\n", highVal)
	fmt.Printf("最低价: %.2f\n", lowVal)
	fmt.Printf("成交量: %d\n", volumeVal)
	fmt.Printf("价格变动: %.2f (%.2f%%)\n", priceChange, percentChange)
	fmt.Printf("日内波动率: %.2f%%\n", volatility)

	// 趋势判断
	if closeVal > openVal {
		fmt.Println("趋势: 上涨")
	} else if closeVal < openVal {
		fmt.Println("趋势: 下跌")
	} else {
		fmt.Println("趋势: 平盘")
	}

	// 波动性评估
	if volatility > 5 {
		fmt.Println("波动性: 高")
	} else if volatility > 2 {
		fmt.Println("波动性: 中")
	} else {
		fmt.Println("波动性: 低")
	}
}

// 辅助函数：字符串转float64
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// 辅助函数：字符串转int64
func parseInt(s string) int64 {
	var i int64
	fmt.Sscanf(s, "%d", &i)
	return i
}
