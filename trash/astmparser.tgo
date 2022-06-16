package astm1384

import (
	"errors"
	"fmt"
)

type CommentedResult struct {
	Result   *Result
	Comments []*Comment
}

func parseAST(token *Token) (*ASTMMessage, error) {

	if token.Type == TokenRoot {
		msg := &ASTMMessage{}
		err := parseHeader(token.Next, msg)
		if err != nil {
			return nil, err
		}
		return msg, nil
	}

	return nil, nil
}

func parseHeader(token *Token, message *ASTMMessage) error {
	if token == nil {
		return errors.New("Premature end of file")
	}

	if token.Type != TokenHeader {
		return errors.New("Header is missing at beginning of input")
	}

	message.Header = token.Data.(*Header)

	for {
		if peekNext(token, TokenPatient) == true {
			var record *Record
			var err error
			token, record, err = parsePatient(token.Next)
			if err != nil {
				return err
			}
			message.Records = append(message.Records, record)
		} else if token == nil {
			return errors.New(fmt.Sprintf("Unexpected end of input"))
		} else if peekNext(token, TokenTerminator) { // the only clean exit here
			token = token.Next
			break
		} else {
			return errors.New(fmt.Sprintf("Unexcpected record type %+v", token.Data))
		}
	}

	return nil
}

func parsePatient(token *Token) (*Token, *Record, error) {
	if token.Type != TokenPatient {
		return token, nil, errors.New(fmt.Sprintf("Expected Patient-Record, but got %+v", token.Data))
	}

	record := &Record{
		Patient:          token.Data.(*Patient),
		OrdersAndResults: make([]*OrderResults, 0),
		Comments:         make([]*Comment, 0)}

	for {
		if peekNext(token, TokenOrder) == true {
			var orderResults *OrderResults
			var err error
			token, orderResults, err = parseOrderResults(token.Next)
			if err != nil {
				return token, nil, err
			}
			record.OrdersAndResults = append(record.OrdersAndResults, orderResults)
		} else if token == nil {
			return token, nil, errors.New(fmt.Sprintf("Unexpected end of input"))
		} else {
			break
		}
	}

	return token, record, nil
}

func parseOrderResults(token *Token) (*Token, *OrderResults, error) {
	if token.Type != TokenOrder {
		return token, nil, errors.New(fmt.Sprintf("Expected Order-Record, but got %+v", token.Data))
	}

	order := &OrderResults{
		Order:    token.Data.(*Order),
		Results:  make([]*CommentedResult, 0),
		Comments: make([]*Comment, 0)}

	for {
		if peekNext(token, TokenResult) == true {
			var commentedResult *CommentedResult
			var err error
			token, commentedResult, err = parseCommentedResult(token.Next)
			if err != nil {
				return token, nil, err
			}
			order.Results = append(order.Results, commentedResult)
		} else if token == nil {
			return token, nil, errors.New(fmt.Sprintf("Unexpected end of input"))
		} else {
			break
		}
	}

	return token, order, nil
}

func parseCommentedResult(token *Token) (*Token, *CommentedResult, error) {
	if token.Type != TokenResult {
		return token, nil, errors.New(fmt.Sprintf("Expected result record (R), but got %+v", token.Data))
	}

	commentedResult := &CommentedResult{
		Result:   token.Data.(*Result),
		Comments: make([]*Comment, 0)}

	for {
		if peekNext(token, TokenComment) == true {
			token = token.Next
			comment := token.Data.(*Comment)
			commentedResult.Comments = append(commentedResult.Comments, comment)
		} else if token == nil {
			return token, nil, errors.New(fmt.Sprintf("Unexpected end of input"))
		} else {
			break
		}
	}

	return token, commentedResult, nil
}

/* See if the next token fits the type */
func peekNext(token *Token, exp TokenType) bool {
	if token.Next == nil {
		return false
	}
	if token.Next.Type == exp {
		return true
	}
	return false
}
