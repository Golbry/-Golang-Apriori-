package main

import (
	"fmt"
)

var curTransaction Transaction
var curCandidates [MaxTxNum]Candidate
var curFrequencies [MaxTxNum]Frequency

func main() {
	initTransaction()
	printTransaction(curTransaction)
	searchFrequencyItemSet()
	printCandidate(1)
	printFrequency(1)
	setIndex := 2
	for {
		aprioriGen(setIndex)
		if curCandidates[setIndex-1].setNum == 0 {
			break
		}
		getFrequencySetIndex(setIndex)
		printCandidate(setIndex)
		printFrequency(setIndex)
		setIndex++
	}
}

func initTransaction() {
	fmt.Println("Please enter the number of records for the current transaction...")
	_, err := fmt.Scanf("%d", &curTransaction.TransactionCount)
	if err != nil {
		panic("fmt.Scanf err!")
	}
	if curTransaction.TransactionCount > MaxTxNum {
		panic("over max transaction limit！")
	}
	i := 0
	for {
		if i >= curTransaction.TransactionCount {
			break
		}
		fmt.Printf("Please enter the number of the %d st/th records\n", i+1)
		_, err = fmt.Scanf("%d", &curTransaction.TransactionItems[i].itemNum)
		if err != nil {
			panic("fmt.Scanf err!")
		}
		fmt.Printf("Please enter Please enter the record item of the %d st/th record (replace with numerical value, order is required)\n", i+1)
		j := 0
		for {
			if j >= curTransaction.TransactionItems[i].itemNum {
				break
			}
			_, err = fmt.Scanf("%d", &curTransaction.TransactionItems[i].items[j])
			if err != nil {
				panic("fmt.Scanf err!")
			}
			j++
		}
		i++
	}
}

func searchFrequencyItemSet() {
	i := 0
	for {
		if i >= curTransaction.TransactionCount { //扫描事务的每条记录
			break
		}
		j := 0
		for {
			//扫描事务的每项记录
			if j >= curTransaction.TransactionItems[i].itemNum {
				break
			}
			candidateItemNum := curCandidates[0].setNum
			isFind := false
			k := 0
			for {
				if k >= candidateItemNum {
					break
				}
				//扫描候选项，检查是否存在
				if curCandidates[0].items[k].items[0] == curTransaction.TransactionItems[i].items[j] {
					curCandidates[0].items[k].itemCount++
					isFind = true
					break
				}
				k++
			}
			//如果没找到，则添加到候选集同时排序
			if !isFind {
				index := 0
				k = 0
				for {
					if k >= candidateItemNum {
						break
					}
					if curCandidates[0].items[k].items[0] > curTransaction.TransactionItems[i].items[j] {
						break
					}
					k++
				}
				index = k
				//之后进行排序
				m := candidateItemNum
				for {
					if m <= index {
						break
					}
					curCandidates[0].items[m] = curCandidates[0].items[m-1]
					m--
				}
				curCandidates[0].items[index].itemCount = 1
				curCandidates[0].items[index].items[0] = curTransaction.TransactionItems[i].items[j]
				curCandidates[0].setNum++
			}
			j++
		}
		i++
	}
	//通过候选集计算频繁集
	setItems := curCandidates[0].setNum
	i = 0
	j := 0
	for {
		if i >= setItems {
			break
		}
		if curCandidates[0].items[i].itemCount >= MinSup {
			curFrequencies[0].items[j] = curCandidates[0].items[i]
			j++
		}
		i++
	}
	curFrequencies[0].setNum = j
}

func aprioriGen(setIndex int) {
	if setIndex < 2 {
		return
	}
	setNum := curFrequencies[setIndex-2].setNum
	candidateNum := 0
	i := 0
	for {
		if i >= setNum {
			break
		}
		j := i + 1
		for {
			if j >= setNum {
				break
			}
			isUnion := true
			k := 0
			for {
				if k >= setIndex-2 {
					break
				}
				if curFrequencies[setIndex-2].items[i].items[k] != curFrequencies[setIndex-2].items[j].items[k] {
					isUnion = false
					break
				}
				k++
			}
			// 若可以链接，按上述链接步骤概念生成
			if isUnion {
				curCandidates[setIndex-1].items[candidateNum] = curFrequencies[setIndex-2].items[i]
				curCandidates[setIndex-1].items[candidateNum].items[setIndex-1] = curFrequencies[setIndex-2].items[j].items[setIndex-2]
				curCandidates[setIndex-1].items[candidateNum].itemCount = 0
				candidateNum++
			}
			j++
		}
		i++
	}
	validCandidateNum := 0
	i = 0
	for {
		if i >= candidateNum {
			break
		}
		curCandidates[setIndex-1].items[validCandidateNum] = curCandidates[setIndex-1].items[i]
		transCount := curTransaction.TransactionCount
		m := 0
		for {
			if m >= transCount {
				break
			}
			itemCount := curTransaction.TransactionItems[m].itemNum
			n := 0
			l := 0
			for {
				if n >= itemCount || l >= setIndex {
					break
				}
				if curCandidates[setIndex-1].items[validCandidateNum].items[l] == curTransaction.TransactionItems[m].items[n] {
					l++
					continue
				} else if curCandidates[setIndex-1].items[validCandidateNum].items[l] < curTransaction.TransactionItems[m].items[n] {
					break
				}
				n++
			}
			if l == setIndex {
				curCandidates[setIndex-1].items[validCandidateNum].itemCount++
			}
			m++
		}
		if curCandidates[setIndex-1].items[validCandidateNum].itemCount > 0 {
			validCandidateNum++
		}
		i++
	}
	curCandidates[setIndex-1].setNum = validCandidateNum
}

func getFrequencySetIndex(setIndex int) {
	i := 0
	j := 0
	setItems := curCandidates[setIndex-1].setNum
	for {
		if i >= setItems {
			break
		}
		if curCandidates[setIndex-1].items[i].itemCount >= MinSup {
			curFrequencies[setIndex-1].items[j] = curCandidates[setIndex-1].items[i]
			j++
		}
		i++
	}
	curFrequencies[setIndex-1].setNum = j
}

func printCandidate(setIndex int) {
	i := 0
	fmt.Printf("candidate set (%v) has %d item sets\n", setIndex, curCandidates[setIndex-1].setNum)
	for {
		if i >= curCandidates[setIndex-1].setNum {
			break
		}
		j := 0
		for {
			if j >= setIndex {
				break
			}
			fmt.Printf("%4d", curCandidates[setIndex-1].items[i].items[j])
			j++
		}
		fmt.Printf("%8d\n", curCandidates[setIndex-1].items[i].itemCount)
		i++
	}
}

func printFrequency(setIndex int) {
	i := 0
	setItems := curFrequencies[setIndex-1].setNum
	fmt.Printf("frequency set (%v) has %d set items\n", setIndex, setItems)
	for {
		if i >= setItems {
			break
		}
		j := 0
		for {
			if j >= setIndex {
				break
			}
			fmt.Printf("%4d", curFrequencies[setIndex-1].items[i].items[j])
			j++
		}
		fmt.Printf("%8d\n", curFrequencies[setIndex-1].items[i].itemCount)
		i++
	}
}

func printTransaction(transaction Transaction) {
	transCount := transaction.TransactionCount
	fmt.Printf("transaction has %d records\n", transCount)
	i := 0
	for {
		if i >= transCount {
			break
		}
		content := fmt.Sprintf("%d :   ", i+1)
		itemCount := transaction.TransactionItems[i].itemNum
		j := 0
		for {
			if j >= itemCount {
				break
			}
			content += fmt.Sprintf("%4d", transaction.TransactionItems[i].items[j])
			j++
		}
		content += "\n"
		fmt.Printf(content)
		i++
	}
}

/*以下是一些定义*/

const MaxTxNum = 200 //事务的最大元素
const MinSup = 2     //最小支持度
// TransactionItem 存储事务的结构
type TransactionItem struct {
	items   [MaxTxNum]int
	itemNum int
}

// Transaction 事务
type Transaction struct {
	TransactionItems [MaxTxNum]TransactionItem
	TransactionCount int
}

// CandidateItem 候选集
type CandidateItem struct {
	items     [MaxTxNum]int
	itemCount int
}

type Candidate struct {
	items  [MaxTxNum]CandidateItem
	setNum int
}

type Frequency struct {
	items  [MaxTxNum]CandidateItem
	setNum int
}
