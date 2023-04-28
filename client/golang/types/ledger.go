package types

type LedgerRawTxInput struct {
	IssueType     string            `index:"0" json:"issueType" xml:"issueType" yaml:"issueType"`              // 签发者：1号流通应用等
	Applicant     string            `index:"5" json:"applicant" xml:"applicant" yaml:"applicant"`              // 申请机构
	ApplicantInst string            `index:"10" json:"applicantInst" xml:"applicantInst" yaml:"applicantInst"` // 申请机构id
	PublisherName string            `index:"15" json:"publisher" xml:"publisherName" yaml:"publisherName"`     // 发布机构
	PublisherInst string            `index:"20" json:"publisherInst" xml:"publisherInst" yaml:"publisherInst"` // 发布机构id
	LedgerNo      string            `index:"25" json:"ledgerNo" xml:"ledgerNo" yaml:"ledgerNo"`                // 存证号
	LedgerDesc    string            `index:"30" json:"ledgerDesc" xml:"ledgerDesc" yaml:"ledgerDesc"`          // 存证名称
	IssueStatus   string            `index:"35" json:"issueStatus" xml:"issueStatus" yaml:"issueStatus"`       // 存证状态
	IssueDate     string            `index:"40" json:"issueDate" xml:"issueDate" yaml:"issueDate"`             // 存证日期
	Content       interface{}       `index:"45" json:"content" xml:"content" yaml:"content"`                   // 存证内容 格式自定义
	Options       map[string]string `index:"50" json:"options" xml:"options" yaml:"options"`                   // 发布存证可选命令
}

type LedgerTxReceipt struct {
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex string `json:"transactionIndex"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	GasUsed          string `json:"gasUsed"`
	ContractAddress  string `json:"contractAddress"`
	Root             string `json:"root"`
	Status           int    `json:"status"`
	From             string `json:"from"`
	To               string `json:"to"`
	Input            string `json:"input"`
	Output           string `json:"output"`
	LogsBloom        string `json:"logsBloom"`
}
