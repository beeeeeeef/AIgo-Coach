package analyze

type AnalyzeRequest struct {
	//网站
	URL string `json:"url" binding:"required"`
	//问题标题
	Title string `json:"title" binding:"required"`
	//问题描述
	Description string `json:"description" binding:"required"`
	Site string `json:"site" binding:"required"`
	//代码语言
	Language string `json:"language" binding:"required"`
	//代码内容
	Code string `json:"code" binding:"required"`
	//模式
	Mode string `json:"mode" binding:"required"`
}


type AnalyzeResponse struct {
	//结果
	Result bool `json:"result"`
	//问题分析
	Data data `json:"data"`
	//错误信息
	Error string `json:"error"`

}

type data struct {
	//做法总结
	Summary string `json:"summary"`
	//问题
	Problem []string `json:"problem"`
	//提示
	Hints []string `json:"hints"`
	//复杂度分析
	Complexity string `json:"complexity"`
}