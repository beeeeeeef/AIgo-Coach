
package analyze
func AnalyzeCode(req AnalyzeRequest) (AnalyzeResponse) {
	 // v0.1 先返回固定 mock 数据
    return AnalyzeResponse{
        Result: true,
        Data: data{
            Summary: "你现在的思路接近暴力枚举...",
            Problem: []string{"..."},
            Hints: []string{"..."},
            Complexity: "...",
        },
    }
}