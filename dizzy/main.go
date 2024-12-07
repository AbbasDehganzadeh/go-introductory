package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

/* A program that solves string manipulation (dizzy) problem
stdin: "b2 B0 o1 !3"
stdout:"Bob!"
*/
const (
	VERSION = "1.0.0"
)

var (
	STRICT  = false
	VERBOSE = true
)

type ByteCipher struct {
	value []byte
	index []byte // TODO: convert to uint
}

func main() {

	cipherInput := []byte{}

	// user input
	buf := bufio.NewScanner(os.Stdin)
	buf.Split(bufio.ScanBytes)
	// bytes_input := new []byte(cipherInput)
	for buf.Scan() {
		cipherInput = append(cipherInput, buf.Bytes()...)
		// fmt.Printf("%T:\t%v", buf.Bytes(), buf.Bytes())
		//
		if buf.Text() == "\n" { // TODO: should quit when no characters are given
			break
		}
	}
	if VERBOSE {
		fmt.Printf("input[Byte]: %v\n", cipherInput)
		fmt.Printf("input[char]: %v\n", string(cipherInput))
	}
	res, err := makeCipher(cipherInput[:len(cipherInput)-1]) // removes mewlines
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	// fmt.Printf("output: %T, \t%v\n", res, res)
	res_str := make([]string, len(cipherInput), len(cipherInput))
	// res_str = append(res_str, strings.Repeat(" ", len(res)))
	for _, chr := range res {
		// fmt.Printf("chr: %v, %d, %d\n", res_str, cap(res_str), len(res_str))
		chr.DecodeByte(&res_str)
	}
	if VERBOSE { // TODO: format output pretty
		fmt.Printf("output: %T, \t%v\n", res_str, res_str)
	} else {
		fmt.Printf("output: \t%v\n", res_str, res_str)
	}
}

func makeCipher(b []byte) ([]ByteCipher, error) {
	// convert []byte to []byte_cipher
	// return error if token is invalid
	bc := make([]ByteCipher, 0, len(b))
	cipher := bytes.Split(b, []byte(" "))
	num_re := regexp.MustCompile(`[\d]+`)
	for i := 0; i < len(cipher); i++ {
		// NOTE: not unnecessary redundancy

		num_idx := num_re.FindIndex(cipher[i])

		if len(num_idx) == 0 && STRICT { // in case of number not found
			return bc, errors.New("NHAVENUM_TOKEN_INVALID: %v must have a number") //, cipher[i])
		} /* else if !STRICT {
			new_bc := ByteCipher{value: cipher[i], index: bc[i-1].index}
			bc = append(bc, new_bc)
			i += 1
			continue
		}*/
		F, L := num_idx[0], num_idx[1]
		if F == 0 && STRICT {
			return bc, errors.New("NHAVECHAR_TOKEN_INVALID: %v must have a character") //, cipher[i])
		}/* else if !STRICT {
			new_bc := ByteCipher{value: bc[i-1].value, index: cipher[i][F:L]}
			bc = append(bc, new_bc)
			i += 1
			continue
		}*/
		// TODO: mutate inputs unless it enables strict mode
		if !STRICT {
		} // TODO: sift down the code above here}
		new_bc := ByteCipher{value: cipher[i][0:F], index: cipher[i][F:L]}
		bc = append(bc, new_bc)

		// fmt.Printf("number index: %T\t%v\n", num_idx, num_idx)
	}

	return bc, nil
}

func (bc *ByteCipher) DecodeByte(s *[]string) error {
	// decodes byteCipher to string
	// returns error if the process is wrong
	var tmpStr string // temporary char string
	tmpStr = string(bc.value)
	idx, _ := strconv.Atoi(string(bc.index))
	(*s)[idx-1] = tmpStr
	//	// idx = idx + 1
	return nil

}
