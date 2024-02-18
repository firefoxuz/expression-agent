package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type MiniTask struct {
	ExpressionId   int    `json:"expression_id"`
	MiniExpression string `json:"task"`
	IsValid        bool   `json:"is_valid"`
	Result         int    `json:"result"`
}

func (miniTask *MiniTask) ResolveMiniTask() (int, error) {
	operands := strings.Split(miniTask.MiniExpression, " ")
	firstValue, _ := strconv.Atoi(operands[0])
	secondValue, _ := strconv.Atoi(operands[1])
	operator := operands[2]

	if operator == "/" && secondValue == 0 {
		miniTask.IsValid = false
		return 0, errors.New("cannot divide by zero")
	}

	result, _ := evaluateOperator(operator, firstValue, secondValue)
	miniTask.Result = result
	miniTask.IsValid = true
	return result, nil
}

func evaluateOperator(oper string, a, b int) (int, error) {
	switch oper {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		return a / b, nil
	case "^":
		return int(math.Pow(float64(a), float64(b))), nil
	default:
		return 0, errors.New("Unknown operator: " + oper)
	}
}

func PublishMiniTaskAnswer(miniTask *MiniTask) {
	rdb, _ := GetConnection()
	channel := fmt.Sprintf("mini_answer_channel_%d", miniTask.ExpressionId)
	content, _ := json.Marshal(miniTask)
	pub := rdb.Publish(context.Background(), channel, string(content))
	fmt.Println(channel, string(content), pub.Err())
}
