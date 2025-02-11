package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ?int=1,2,3  // OR
// ?int_arr=1,2,3   // all arrays with any of these values
// ?str=1&str1=2&str1=3  // OR
// ?str_arr=1&str_arr=2&str_arr=3 // all arrays with any of these values

func getQueryClauseArgs(params url.Values, columnMap map[string]Column, tableName string) (string, []any, error) {
	var argNum = 1
	var conditions []string
	var args []any
	var clause string
	var err error

	// Conditions
	for columnName, column := range columnMap {
		if column.Hash {
			continue
		}

		queryArr := params[columnName]

		if len(queryArr) == 0 {
			continue
		}

		// single key if not string or string arr
		if !strings.HasSuffix(column.DataType, "String") {
			queryArr = strings.Split(queryArr[0], ",")
		}

		if args, err = appendArgs(args, queryArr, column.DataType); err != nil {
			errorMessage := fmt.Sprintf(`error while parsing %s args: %v`, columnName, err)
			return clause, args, errors.New(errorMessage)
		}

		argPlaceHolders := getArgPlaceHolders(argNum, len(queryArr))
		argNum += len(queryArr)

		var condition string
		if strings.HasPrefix(column.DataType, "[]") {
			condition = fmt.Sprintf(`"%s"."%s" && ARRAY[%v]::%s`, tableName, columnName, argPlaceHolders, column.pgType)
		} else {
			condition = fmt.Sprintf(`"%s"."%s" IN (%v)`, tableName, columnName, argPlaceHolders)
		}
		conditions = append(conditions, condition)
	}

	if len(conditions) > 0 {
		clause += fmt.Sprintf(` WHERE %s`, strings.Join(conditions, " AND "))
	}

	// Order By
	orderVal := strings.Split(params.Get("__order"), ",")
	orderBy := []string{}

	for _, columnName := range orderVal {
		order := "ASC"

		if strings.HasPrefix(columnName, "-") {
			columnName = columnName[1:]
			order = "DESC"
		}

		if _, ok := columnMap[columnName]; ok {
			orderBy = append(orderBy, fmt.Sprintf(`"%s"."%s" %s`, tableName, columnName, order))
		}
	}

	if len(orderBy) > 0 {
		clause += fmt.Sprintf(` ORDER BY %s`, strings.Join(orderBy, ", "))
	}

	skipCount, _ := strconv.Atoi(params.Get("__skip"))
	args = append(args, skipCount)
	clause += fmt.Sprintf(` OFFSET $%d ROWS`, len(args))

	limitCount, _ := strconv.Atoi(params.Get("__limit"))
	args = append(args, limitCount+1)
	clause += fmt.Sprintf(` FETCH FIRST $%d ROWS ONLY`, len(args))

	return clause, args, nil
}

func appendArgs(argsList []any, values []string, datatype string) ([]any, error) {
	isInt := strings.HasSuffix(datatype, "Int")
	isFloat := strings.HasSuffix(datatype, "Float")

	for _, value := range values {
		if isInt {
			parsed, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return argsList, err
			}
			argsList = append(argsList, parsed)
		} else if isFloat {
			parsed, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return argsList, err
			}
			argsList = append(argsList, parsed)
		} else {
			argsList = append(argsList, value)
		}
	}

	return argsList, nil
}

func getArgPlaceHolders(start, count int) string {
	placeholders := []string{}

	for i := 1; i <= count; i++ {
		placeholders = append(placeholders, fmt.Sprintf(`$%d`, start))
		start++
	}

	return strings.Join(placeholders, ", ")
}

// ?id=1 OR ?id=text OR ?id=6.4 OR ?id=3.5  // non-array variables
// ?id=1,2,3   // non-string array
// ?id=text1&id=text2&id=text3 // string array

func getPkParam(params url.Values, datatype string) string {
	if datatype == "[]CustomNullString" {
		values := params["id"]
		return fmt.Sprintf(`{%s}`, strings.Join(values, ","))
	}

	id := params.Get("id")

	if strings.HasPrefix(datatype, "[]") {
		id = "{" + id + "}"
	}

	return id
}
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func comparePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func hashData(strings []*CustomNullString, stringArrays [][]*CustomNullString) error {
	for _, string_ := range strings {
		hashed, err := hashPassword(string_.String)
		if err != nil {
			return err
		}
		string_.String = hashed
	}

	for _, subArr := range stringArrays {
		if err := hashData(subArr, nil); err != nil {
			return err
		}
	}

	return nil
}

func getSignedToken(username, role, college_id, course_id, branch_id string) (string, error) {

	claims := CustomJwtClaims{
		username, role, college_id, course_id, branch_id, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	return signedString, err
}

func validateSignedToken(signedToken string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("failed to parse claims")
}

func authorizeRequest(r *http.Request, allowedRoles []string) (jwt.MapClaims, error) {
	emptyClaim := jwt.MapClaims{
		"username": "",

		"role":       "",
		"college_id": "", "course_id": "", "branch_id": ""}

	cookie, err := r.Cookie("access_token")
	if err != nil {
		return emptyClaim, err
	}

	claims, err := validateSignedToken(cookie.Value)
	if err != nil {
		return emptyClaim, err
	}

	role := claims["role"].(string)

	if len(allowedRoles) > 0 && !slices.Contains(allowedRoles, role) {
		return emptyClaim, fmt.Errorf("%s role is not authorized", role)
	}

	return claims, nil
}

func validateProtectedField[K comparable](info map[K][]string, value K, role string) bool {
	allowedRoles, ok := info[value]

	if !ok {
		return true
	}

	return slices.Contains(allowedRoles, role)
}

func validateArrProtectedField[K comparable](info map[K][]string, arr []K, role string) bool {
	for _, item := range arr {
		allowedRoles, ok := info[item]

		if ok && !slices.Contains(allowedRoles, role) {
			return false
		}
	}

	return true
}

func protectedFieldsQueryArgs[K comparable](info map[K][]string, role string) []any {
	var args []any

	for value, allowedRoles := range info {
		if !slices.Contains(allowedRoles, role) {
			args = append(args, value)
		}
	}

	return args
}
